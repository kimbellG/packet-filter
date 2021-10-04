#include <linux/stddef.h>
#include <linux/swab.h>

#include <uapi/linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>

#if __BYTE_ORDER__ == __ORDER_LITTLE_ENDIAN__
#define _bpf_ntohs(x) __builtion_bswap16(x)
#define __bpf_htons(x) __builtin_bswap16(x)

#define _bpf_constant_ntohs(x) ___constant_swab16(x)
# define __bpf_constant_htons(x) ___constant_swab16(x)

#elif __BYTE_ORDER__ == __ORDER_BIG_ENDIAN__
#define _bpf_ntohs(x) (x)
#define _bpf_constant_ntohs(x)
#else
#error "define __BYTE_ORDER__"
#endif

#define bpf_htons(x)	\
	(__builtin_constant_p(x) ?	\
	 _bpf_constant_ntohs(x) : _bpf_ntohs(x))

#define INTERNAL static __attribute__((always_inline))


BPF_HASH(pacinfo, u8, u64, 1);
BPF_HASH(protocol_blacklist, u16, u16, 256); 


INTERNAL void count_increment()
{
	u8 count = COUNTKEY;
	u64 zero = 0;
	u64 *val = pacinfo.lookup_or_try_init(&count, &zero);
	if (val) {
		lock_xadd(val, 1);
	}

}

INTERNAL int process_ip(struct iphdr *ip, void *data_end)
{
	bpf_trace_printk("received ip packcage. ID of next protocol: 0x%x, icmp: 0x%x", ip->protocol);

	u16 id = ip->protocol;
	u16 *val = protocol_blacklist.lookup(&id);
	if (val)
	{
		bpf_trace_printk("drop package");
		return XDP_DROP;
	}

	return XDP_PASS;
}

INTERNAL int process_ether(struct ethhdr *ether, void *data_end)
{
	bpf_trace_printk("received ethernet package. ID of next proto: 0x%x", ether->h_proto);

	if (ether->h_proto != bpf_ntohs(ETH_P_IP))
	{
		return XDP_PASS;
	}

	struct iphdr *ip = (struct iphdr *)(ether + 1);
	if ((void *)(ip + 1) > data_end)
	{
		return XDP_PASS;
	}

	return process_ip(ip, data_end);
}


int filter_main(struct xdp_md *ctx)
{
	count_increment();
	
	struct ethhdr *ether = (struct ethhdr*)(void*)ctx->data;
	if ((void *)(ether + 1) > (void *)ctx->data_end)
	{
		return XDP_PASS;
	}

	return process_ether(ether, (void *)ctx->data_end);
}
