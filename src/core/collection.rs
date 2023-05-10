use super::index::Index;

pub struct Collection {
    pub name: String,
    pub index: Box<dyn Index>,
}

impl Collection {
    pub fn new(name: &str, index: Box<dyn Index>) -> Collection {
        Collection {
            name: name.to_string(),
            index,
        }
    }
}
