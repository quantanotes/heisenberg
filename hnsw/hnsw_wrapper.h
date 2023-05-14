#include <stdbool.h>
#ifdef __cplusplus
extern "C"
{
#endif
    typedef void *HNSW;
    HNSW initHNSW(int dim, unsigned long int max_elements, int M, int ef_construction, int rand_seed, int stype);
    HNSW loadHNSW(char *location, int dim, int stype);
    bool saveHNSW(HNSW index, char *location);
    void freeHNSW(HNSW index);
    void addPoint(HNSW index, float *vec, unsigned long int id);
    void deletePoint(HNSW index, unsigned long int id);
    int search(HNSW index, float *vec, int k, unsigned long int *ids, float *dist);
#ifdef __cplusplus
}
#endif