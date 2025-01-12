CREATE TABLE node_closure (
    ancestor UUID NOT NULL,
    descendant UUID NOT NULL,
    depth INT NOT NULL,
    PRIMARY KEY (ancestor, descendant),
    FOREIGN KEY (ancestor) REFERENCES nodes(id),
    FOREIGN KEY (descendant) REFERENCES nodes(id)
);
