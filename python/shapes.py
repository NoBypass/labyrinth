from enum import Enum

from PIL import ImageDraw

from python.consts import offset, multiplier


class Direction(Enum):
    UP = (-1, 0)
    DOWN = (1, 0)
    LEFT = (0, -1)
    RIGHT = (0, 1)

class Line:
    def __init__(self, field: 'Field', d: Direction, color='black'):
        self.x1 = field.x
        self.y1 = field.y
        self.x2 = field.x
        self.y2 = field.y

        if d == Direction.UP:
            self.y1 -= 1
            self.x2 += 1
            self.y2 -= 1
        elif d == Direction.DOWN:
            self.x2 += 1
        elif d == Direction.LEFT:
            self.y2 -= 1
        elif d == Direction.RIGHT:
            self.x1 += 1
            self.x2 += 1
            self.y2 -= 1

        self.y1 += 1
        self.y2 += 1

        self.x1 *= multiplier
        self.y1 *= multiplier
        self.x2 *= multiplier
        self.y2 *= multiplier

        self.x1 += offset
        self.y1 += offset
        self.x2 += offset
        self.y2 += offset

        self.shown = True
        self.color = color

    def hide(self):
        self.shown = False

    def draw(self, d: ImageDraw):
        if self.shown:
            d.line([(self.x1, self.y1), (self.x2, self.y2)], fill=self.color)


class Field:
    def __init__(self, x, y):
        self.border: dict[Direction, Line] = {}
        self.x = x
        self.y = y

    def set_border(self, direction: Direction, border: Line):
        self.border[direction] = border

    def move(self, grid: 'Grid', direction: Direction):
        m: tuple[int, int] = direction.value
        new_x, new_y = self.x + m[1], self.y + m[0]
        if 0 <= new_x < grid.size and 0 <= new_y < grid.size:
            return grid[new_y][new_x]
        return None

class Grid:
    def __init__(self, size: int):
        self.grid = [[Field(x, y) for x in range(size)] for y in range(size)]
        self.size = size
        self.lines = {}

        for y in range(size):
            for x in range(size):
                f = self.grid[y][x]
                f.border = {
                    Direction.RIGHT: Line(f, Direction.RIGHT),
                    Direction.LEFT: Line(f, Direction.LEFT),
                    Direction.DOWN: Line(f, Direction.DOWN),
                    Direction.UP: Line(f, Direction.UP),
                }
                if y != 0:
                    f.border[Direction.UP] = self.grid[y - 1][x].border[Direction.DOWN]
                if x != 0:
                    f.border[Direction.LEFT] = self.grid[y][x - 1].border[Direction.RIGHT]
                for l in f.border.values():
                    self.lines[l] = None
                self.grid[y][x] = f

        self.lines = list(self.lines.keys())

    def __getitem__(self, item):
        return self.grid[item]