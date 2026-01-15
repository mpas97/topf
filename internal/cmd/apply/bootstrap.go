package apply

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/postfinance/topf/internal/topf"
	"github.com/postfinance/topf/pkg/config"
	"github.com/siderolabs/go-retry/retry"
	"github.com/siderolabs/talos/pkg/machinery/api/machine"
)

// bootstrap initiates the ETCD bootstrap process and waits for nodes to stabilize
func bootstrap(ctx context.Context, logger *slog.Logger, nodes []*topf.Node) error {
	if len(nodes) == 0 || nodes[0].Node.Role != config.RoleControlPlane {
		return errors.New("bootstrap requires at least 1 control plane node")
	}

	logger.Info("starting bootstrap process", "timeout", "10 minutes")

	err := retry.Constant(time.Minute*10, retry.WithErrorLogging(logger.Enabled(ctx, slog.LevelDebug))).RetryWithContext(ctx, func(ctx context.Context) error {
		// bootstrap needs to be called on any CP node, we take the first one
		nodeClient, err := nodes[0].Client(ctx)
		if err != nil {
			return retry.ExpectedErrorf("couldn't get client for bootstrap: %w", err)
		}
		defer nodeClient.Close()

		// Check if etcd is already bootstrapped
		membersResp, err := nodeClient.MachineClient.EtcdMemberList(ctx, &machine.EtcdMemberListRequest{})
		if err == nil && len(membersResp.GetMessages()) > 0 {
			// Check if any message contains etcd members
			for _, msg := range membersResp.GetMessages() {
				if len(msg.GetMembers()) > 0 {
					logger.Info("etcd already bootstrapped", "member_count", len(msg.GetMembers()))
					return nil // Already bootstrapped - success
				}
			}
		}

		// Not bootstrapped or error checking - attempt bootstrap
		_, err = nodeClient.MachineClient.Bootstrap(ctx, &machine.BootstrapRequest{})
		if err != nil {
			return retry.ExpectedError(err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("bootstrap failed: %w", err)
	}

	logger.Info("etcd bootstrap completed successfully")

	return nil
}
