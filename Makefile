#!/usr/bin/make --no-print-directory --jobs=1 --environment-overrides -f

CORELIB_PKG := go-corelibs/x-sync
VERSION_TAGS += MAIN
MAIN_MK_SUMMARY := ${CORELIB_PKG}
MAIN_MK_VERSION := v0.1.0

DEPS += golang.org/x/perf/cmd/benchstat

STATS_BENCH       := testdata/bench
STATS_FILE        := ${STATS_BENCH}/${MAIN_MK_VERSION}
STATS_PATH        := ${STATS_BENCH}/${MAIN_MK_VERSION}-d
STATS_FILE_OUTPUT := ${STATS_BENCH}/${MAIN_MK_VERSION}-d/output
STATS_FILE_GOLANG := ${STATS_BENCH}/${MAIN_MK_VERSION}-d/golang
STATS_FILE_X_SYNC := ${STATS_BENCH}/${MAIN_MK_VERSION}-d/x-sync

STATS_FILES += ${STATS_FILE}
STATS_FILES += ${STATS_FILE_OUTPUT}
STATS_FILES += ${STATS_FILE_GOLANG}
STATS_FILES += ${STATS_FILE_X_SYNC}

.PHONE += benchmark
.PHONY += benchstats-history
.PHONY += benchstats-versus

define _perl_benchmark_filter
if (m!_Golang!) { \
  s/_Golang//; \
  print STDERR "$$_"; \
} elsif (m!_X_Sync!) { \
	s/_X_Sync//; \
  print STDOUT "$$_"; \
} else { \
  print STDOUT "$$_"; \
  print STDERR "$$_"; \
}; $$_="";
endef

include CoreLibs.mk

benchmark: export BENCH_COUNT=500
benchmark:
	@rm -fv    "${STATS_FILE}" || true
	@rm -rfv   "${STATS_PATH}" || true
	@mkdir -vp "${STATS_PATH}"
	@$(MAKE) bench | egrep -v '^make' > "${STATS_FILE_OUTPUT}"
	@cat "${STATS_FILE_OUTPUT}" \
			| grep -v "_Golang" \
			> "${STATS_FILE}"
	@cat "${STATS_FILE_OUTPUT}" \
			| perl -pe '$(call _perl_benchmark_filter)' \
			> "${STATS_FILE_X_SYNC}" \
			2> "${STATS_FILE_GOLANG}"
	@shasum ${STATS_FILES}

benchstats-history:
	@pushd ${STATS_BENCH} > /dev/null \
		&& ${CMD} benchstat \
			`ls | egrep -v '\-d$$' | sort -V` \
		&& popd > /dev/null

benchstats-versus:
	@pushd ${STATS_PATH} > /dev/null \
		&& ${CMD} benchstat \
			`basename ${STATS_FILE_GOLANG}` \
			`basename ${STATS_FILE_X_SYNC}` \
		&& popd > /dev/null
