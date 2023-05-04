#include <iostream>
#include "hnswlib/hnswlib.h"
#include "hnsw_wrapper.h"
#include <thread>
#include <atomic>

hnswlib::SpaceInterface<float> getSpace(char *space_type)
{
    switch (space_type)
    {
    case "ip":
        return new hnswlib::InnerProductSpace(dim);
        break;
    case "cosine":
        return new hnswlib::InnerProductSpace(dim);
    case "l2":
        return new hnswlib::L2Space(dim);
    default:
        return;
    }
}

HNSW init(int dim, unsigned long int max, int m, int ef, int seed, char space_type)
{
    hnswlib::SpaceInterface<float> *space = getSpace(space_type);
    hnswlib::HierarchicalNSW<float> *hnsw = new hnswlib::HierarchicalNSW<float>(space, max, m, ef, seed, true);
    return (void *)hnsw;
}

HNSW load(char *path, int dim, char space_type)
{
    hnswlib::SpaceInterface<float> *space = getSpace(space_type);
    hnswlib::HierarchicalNSW<float> *hnsw = new hnswlib::HierarchicalNSW<float>(space, std::string(location), false, 0, true);
    return (void *)hnsw;
}

bool save(HNSW hnsw, char *path)
{
    try
    {
        ((hnswlib::HierarchicalNSW<float> *)hnsw)->saveIndex(std::string(path));
        return true;
    }
    catch (const std::exception_ptr e)
    {
        return false;
    }
}

void free(HNSW hnsw)
{
    delete (hnswlib::HierarchicalNSW<float> *)hnsw;
}

void add(HNSW hnsw, float *vec, unsigned long int id, float dist)
{
    ((hnswlib::HierarchicalNSW<float> *)hnsw)->addPoint(vec, id, true);
}

void remove(HNSW hnsw, unsigned long int id)
{
    ((hnswlib::HierarchicalNSW<float> *)hnsw)->markDelete(id);
}

void search(HNSW hnsw, k int)
{
    std::priority_queue<std::pair<float, hnswlib::labeltype>> result = ((hnswlib::HierarchicalNSW<float> *)hnsw)->searchKnn(vec, k);
    int size = result.size();
    float ids[size];
    for (int i = 0; i < size; i++)
    {
        ids[i] = result.top().second;
        result.pop();
    }
    return ids;
}