LEFT = (0, -1)
RIGHT = (0, 1)
UP = (-1, 0)
DOWN = (1, 0)
UNIT_VECTORS = (LEFT, RIGHT, UP, DOWN)


def invert_uv(v):
    return (v[0] * -1, v[1] * -1)


def add_uv(j, i, v):
    return j + v[0], i + v[1]
