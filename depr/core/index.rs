pub enum IndexError {
    DimensionError(usize, usize),
}

pub trait Index {
    fn new(dim: usize) -> Result<Self, IndexError>
    where
        Self: Sized;
    fn insert(&mut self, key: String, vector: Vec<f32>) -> Result<(), IndexError>;
    fn remove(&mut self, key: String) -> Result<(), IndexError>;
    fn search(&self, query: Vec<f32>, k: usize) -> Result<Vec<&str>, IndexError>;
}

impl Index for Box<dyn Index> {
    fn new(dim: usize) -> Result<Self, IndexError>
    where
        Self: Sized,
    {
        Ok(Box::new(Self::new(dim)?))
    }

    fn insert(&mut self, key: String, vector: Vec<f32>) -> Result<(), IndexError> {
        self.as_mut().insert(key, vector)
    }

    fn remove(&mut self, key: String) -> Result<(), IndexError> {
        self.as_mut().remove(key)
    }

    fn search(&self, query: Vec<f32>, k: usize) -> Result<Vec<&str>, IndexError> {
        self.as_ref().search(query, k)
    }
}