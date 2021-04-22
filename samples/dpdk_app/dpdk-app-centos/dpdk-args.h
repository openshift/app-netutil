// SPDX-License-Identifier: Apache-2.0
// Copyright(c) 2021 Red Hat, Inc.

#ifndef __DPDK_ARGS_H__
#define __DPDK_ARGS_H__

typedef enum {
	DPDK_APP_TESTPMD = 1,
	DPDK_APP_L3FWD,
	DPDK_APP_L2FWD,
	DPDK_APP_OTHER
} eDpdkAppType;

extern char** GetArgs(int *pArgc, eDpdkAppType appType);

#endif  /* __DPDK_ARGS_H__ */
