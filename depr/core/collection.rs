use super::index::Index;

pub struct Collection<I: Index> {
    name: String,
    pub index: I,
}

impl<I: Index> Collection<I> {
    pub fn new(name: String, index: I) -> Collection<I> {
        Collection { name, index }
    }

    pub fn name(&self) -> &str {
        &self.name
    }
}
