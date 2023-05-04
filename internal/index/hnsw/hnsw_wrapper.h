#ifdef __cplusplus
extern "C"
{
#endif
    typedef void *HNSW;
    HNSW init(int dim, unsigned long int max_elements, int M, int ef_construction, int rand_seed, char stype);
    HNSW load(char *location, int dim, char stype);
    HNSW save(HNSW index, char *location);
    void free(HNSW index);
    void add(HNSW index, float *vec, unsigned long int label);
    void remove(HNSW index, float *vec, unsigned long int label);
    int search(HNSW index, float *vec, int N, unsigned long int *label, float *dist);
#ifdef __cplusplus
}
#endif