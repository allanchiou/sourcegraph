package main

func PreciseCodeIntelIndexer() *Container {
	return &Container{
		Name:        "precise-code-intel-indexer",
		Title:       "Precise Code Intel Indexer",
		Description: "Automatically indexes from popular, active Go repositories.",
		Groups: []Group{
			{
				Title: "General",
				Rows: []Row{
					{
						{
							Name:              "index_queue_size",
							Description:       "index queue size",
							Query:             `max(src_index_queue_indexes_total)`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 100},
							PanelOptions:      PanelOptions().LegendFormat("indexes queued for processing"),
							PossibleSolutions: "none",
						},
						{
							Name:              "index_queue_growth_rate",
							Description:       "index queue growth rate every 5m",
							Query:             `sum(increase(src_index_queue_indexes_total[30m])) / sum(increase(src_index_queue_processor_total[30m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 5},
							PanelOptions:      PanelOptions().LegendFormat("index queue growth rate"),
							PossibleSolutions: "none",
						},
						{
							Name:        "index_process_errors",
							Description: "index process errors every 5m",
							// TODO(efritz) - ensure these differentiate unexpected repo layout and system errors
							Query:             `sum(increase(src_index_queue_processor_errors_total[5m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("errors"),
							PossibleSolutions: "none",
						},
					},
					{
						{
							Name:        "99th_percentile_store_duration",
							Description: "99th percentile successful database query duration over 5m",
							// TODO(efritz) - ensure these exclude error durations
							Query:             `histogram_quantile(0.99, sum by (le)(rate(src_code_intel_store_duration_seconds_bucket{job="precise-code-intel-indexer"}[5m])))`,
							DataMayNotExist:   true,
							DataMayBeNaN:      true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("store operation").Unit(Seconds),
							PossibleSolutions: "none",
						},
						{
							Name:              "store_errors",
							Description:       "database errors every 5m",
							Query:             `increase(src_code_intel_store_errors_total{job="precise-code-intel-indexer"}[5m])`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("store operation"),
							PossibleSolutions: "none",
						},
					},
				},
			},
			{
				Title:  "Indexability updater and index scheduler",
				Hidden: true,
				Rows: []Row{
					{
						{
							Name:              "indexability_updater_errors",
							Description:       "indexability updater errors every 5m",
							Query:             `sum(increase(src_indexability_updater_errors_total[5m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("errors"),
							PossibleSolutions: "none",
						},
						{
							Name:              "index_scheduler_errors",
							Description:       "index scheduler errors every 5m",
							Query:             `sum(increase(src_index_scheduler_errors_total[5m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("errors"),
							PossibleSolutions: "none",
						},
					},
				},
			},
			{
				Title:  "Index resetter - re-queues indexes that did not complete processing",
				Hidden: true,
				Rows: []Row{
					{
						{
							Name:              "processing_indexes_reset",
							Description:       "indexes reset to queued state every 5m",
							Query:             `sum(increase(src_index_queue_resets_total[5m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("indexes"),
							PossibleSolutions: "none",
						},
						{
							Name:              "processing_indexes_reset_failures",
							Description:       "indexes errored after repeated resets every 5m",
							Query:             `sum(increase(src_index_queue_max_resets_total[5m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("indexes"),
							PossibleSolutions: "none",
						},
						{
							Name:              "index_resetter_errors",
							Description:       "index resetter errors every 5m",
							Query:             `sum(increase(src_index_queue_reset_errors_total[5m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("errors"),
							PossibleSolutions: "none",
						},
					},
				},
			},
			{
				Title:  "Janitor - synchronizes database and filesystem and keeps free space on disk",
				Hidden: true,
				Rows: []Row{
					{
						{
							Name:              "janitor_errors",
							Description:       "janitor errors every 5m",
							Query:             `sum(increase(src_indexer_janitor_errors_total[5m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("errors"),
							PossibleSolutions: "none",
						},
						{
							Name:              "janitor_indexes_removed",
							Description:       "index records removed every 5m",
							Query:             `sum(increase(src_indexer_janitor_index_records_removed_total[5m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("records removed"),
							PossibleSolutions: "none",
						},
					},
				},
			},
			{
				Title:  "Internal service requests",
				Hidden: true,
				Rows: []Row{
					{
						{
							Name:              "99th_percentile_gitserver_duration",
							Description:       "99th percentile successful gitserver query duration over 5m",
							Query:             `histogram_quantile(0.99, sum by (le,category)(rate(src_gitserver_request_duration_seconds_bucket{job="precise-code-intel-indexer"}[5m])))`,
							DataMayNotExist:   true,
							DataMayBeNaN:      true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("{{category}}").Unit(Seconds),
							PossibleSolutions: "none",
						},
						{
							Name:              "gitserver_error_responses",
							Description:       "gitserver error responses every 5m",
							Query:             `sum by (category)(increase(src_gitserver_request_duration_seconds_count{job="precise-code-intel-indexer",code!~"2.."}[5m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 5},
							PanelOptions:      PanelOptions().LegendFormat("{{category}}"),
							PossibleSolutions: "none",
						},
					},
					{
						sharedFrontendInternalAPIErrorResponses("precise-code-intel-indexer"),
					},
				},
			},
			{
				Title:  "Container monitoring (not available on server)",
				Hidden: true,
				Rows: []Row{
					{
						sharedContainerRestarts("precise-code-intel-indexer"),
						sharedContainerMemoryUsage("precise-code-intel-indexer"),
						sharedContainerCPUUsage("precise-code-intel-indexer"),
					},
				},
			},
			{
				Title:  "Provisioning indicators (not available on server)",
				Hidden: true,
				Rows: []Row{
					{
						sharedProvisioningCPUUsage7d("precise-code-intel-indexer"),
						sharedProvisioningMemoryUsage7d("precise-code-intel-indexer"),
					},
					{
						sharedProvisioningCPUUsage5m("precise-code-intel-indexer"),
						sharedProvisioningMemoryUsage5m("precise-code-intel-indexer"),
					},
				},
			},
		},
	}
}
