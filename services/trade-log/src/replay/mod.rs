use serde::Serialize;
#[derive(Debug, Serialize)]
pub struct ReplayReport {
    pub query_id: String,
    pub old_parser_version: String,
    pub parser_version: String,
    pub source_records: usize,
    pub fact_records: usize,
    pub comparison: String,
    pub differences: Vec<String>,
}
