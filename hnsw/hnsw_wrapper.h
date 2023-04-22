// hnsw_wrapper.h
#ifdef __cplusplus
extern "C"
{
#endif
    typedef void *HNSW;
    HNSW initHNSW(int dim, unsigned long int max_elements, int M, int ef_construction, int rand_seed, char stype);
    HNSW loadHNSW(char *location, int dim, char stype);
    HNSW saveHNSW(HNSW index, char *location);
    void freeHNSW(HNSW index);
    void addPoint(HNSW index, float *vec, unsigned long int label);
    int searchKnn(HNSW index, float *vec, int N, unsigned long int *label, float *dist);
    void setEf(HNSW index, int ef);
#ifdef __cplusplus
}
#endif