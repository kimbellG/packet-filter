#include <uapi/linux/bpf.h>
#include <linux/if_ether.h>

BPF_HASH(MAPNAME, u8, u64, 1);


int FUNCNAME(struct xdp_md *ctx)
{
	u8 count = COUNTKEY;
	u64 zero = 0;
	u64 *val = MAPNAME.lookup_or_try_init(&count, &zero);
	if (val) {
		lock_xadd(val, 1);
	}

	return XDP_PASS;
}
