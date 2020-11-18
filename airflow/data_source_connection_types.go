package airflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const schemaDataSourceConnectionTypes = "airflow_connection_types"

const (
	mkDataSourceConnectionTypesTypes = "types"
)

func dataSourceConnectionTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectionTypesRead,
		Schema: map[string]*schema.Schema{
			mkDataSourceConnectionTypesTypes: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceConnectionTypesRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics

	_ = d.Set(mkDataSourceConnectionTypesTypes, getConnectionTypesAsList())

	d.SetId("airflow")

	return diags
}

func getConnectionTypesAsList() []string {
	return []string{
		"aws",
		"azure_batch",
		"azure_container_instances",
		"azure_cosmos",
		"azure_data_explorer",
		"azure_data_lake",
		"azure",
		"cassandra",
		"cloudant",
		"databricks",
		"docker",
		"elasticsearch",
		"emr",
		"exasol",
		"facebook_social",
		"fs",
		"ftp",
		"gcpcloudsql",
		"google_cloud_platform",
		"grpc",
		"hdfs",
		"hive_cli",
		"hive_metastore",
		"hiveserver2",
		"http",
		"imap",
		"jdbc",
		"jenkins",
		"jira",
		"kubernetes",
		"livy",
		"mesos_framework-id",
		"mongo",
		"mssql",
		"mysql",
		"odbc",
		"oracle",
		"pig_cli",
		"postgres",
		"presto",
		"qubole",
		"redis",
		"s3",
		"samba",
		"segment",
		"snowflake",
		"spark",
		"sqlite",
		"sqoop",
		"ssh",
		"tableau",
		"vault",
		"vertica",
		"wasb",
		"yandexcloud",
	}
}
