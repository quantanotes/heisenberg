use hora::{
    core::{ann_index::ANNIndex, node::Node},
    index::{hnsw_idx::HNSWIndex, hnsw_params::HNSWParams},
};

use crate::core::{model::*, index::Index};

pub type HoraHNSWIndexOptions = HNSWParams<f32>;

pub struct HoraHNSWIndex {
    dimension: usize,
    index: HNSWIndex<f32, String>,
}

impl Index for HoraHNSWIndex {
    fn new(dimension: usize) -> Result<HoraHNSWIndex, Error> {
        Ok(HoraHNSWIndex {
            dimension,
            index: hora::index::hnsw_idx::HNSWIndex::new(
                dimension,
                &HNSWParams::default().has_deletion(true),
            ),
        })
    }

    fn insert(&mut self, idx: String, vector: Vec<f32>) -> Result<(), Error> {
        Ok(())


        // self.index
        //     .add_node(&Node {
        //         vectors: vec![vector],
        //         idx: Some(idx),
        //     })
        //     .map_err(Into::into)
    }

    fn remove(&mut self, idx: String) -> Result<(), Error> {
        Ok(())//self.index.(idx)
    }

    fn search(&self, query: Vec<f32>, k: usize) -> Result<Vec<String>, Error> {
        // if query.len() != self.dimension {
        //     Err(IndexError::DimensionError(query.len(), self.dimension))
        // }

        // Ok(self.index.search(&query, k))
        Ok(Vec::new())
    }
}
