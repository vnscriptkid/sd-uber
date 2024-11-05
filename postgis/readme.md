# GiST (Generalized Search Tree)
- R-tree, multi-dimensional objects
- How it works?
  - Bounding box for geometries
  - Tree structure
  - Efficient search: due to ability to skip sub-trees (eliminate large portions of the tree that don't match)

# SP-GiST (Space-Partitioned GiST)
- Non-overlapping regions
- Data structure: quad-tree, k-d tree, prefix-tree