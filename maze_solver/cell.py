from config import DEFAULTCELLBORDERCOLOR
import constants as const


class Point:
    def __init__(self, x, y):
        self._x = x
        self._y = y

    def __repr__(self):
        return f"Point({self._x}, {self._y})"


class Line:
    def __init__(self, start, end):
        self._start = start
        self._end = end

    def draw(self, canvas, fill_color):
        canvas.create_line(
            self._start._x,
            self._start._y,
            self._end._x,
            self._end._y,
            fill=fill_color,
            width=2,
        )
        canvas.pack()


class Cell:
    def __init__(self, top_left_x, top_left_y, x_size, y_size, win=None, visited=False):
        self.has_left_wall = True
        self.has_right_wall = True
        self.has_top_wall = True
        self.has_bottom_wall = True

        self._top_left = Point(top_left_x, top_left_y)
        self._top_right = Point(top_left_x + x_size, top_left_y)
        self._bot_left = Point(top_left_x, top_left_y + y_size)
        self._bot_right = Point(top_left_x + x_size, top_left_y + y_size)
        self._center = Point((2*top_left_x + x_size) / 2,
                             (2*top_left_y + y_size) / 2)

        self._win = win
        self.visited = visited

    def draw(self):
        l_color, r_color, t_color, b_color = "light gray", "light gray", "light gray", "light gray"

        if self.has_left_wall:
            l_color = DEFAULTCELLBORDERCOLOR
        if self.has_right_wall:
            r_color = DEFAULTCELLBORDERCOLOR
        if self.has_top_wall:
            t_color = DEFAULTCELLBORDERCOLOR
        if self.has_bottom_wall:
            b_color = DEFAULTCELLBORDERCOLOR

        self._win.draw_line(Line(self._top_left, self._bot_left), l_color)
        self._win.draw_line(Line(self._top_right, self._bot_right), r_color)
        self._win.draw_line(Line(self._top_left, self._top_right), t_color)
        self._win.draw_line(Line(self._bot_left, self._bot_right), b_color)

    def draw_move(self, to_cell, undo=False):
        if undo:
            color = "grey"
        else:
            color = "red"
        self._win.draw_line(Line(self._center, to_cell._center), color)

    def has_visited(self):
        self.visited = True

    def knock_down_wall(self, unit_vector):
        if unit_vector == const.LEFT:
            self.has_left_wall = False
        if unit_vector == const.RIGHT:
            self.has_right_wall = False
        if unit_vector == const.UP:
            self.has_top_wall = False
        if unit_vector == const.DOWN:
            self.has_bottom_wall = False

    def is_dir_clear(self, unit_vector):
        if unit_vector == const.LEFT and not self.has_left_wall:
            return True
        if unit_vector == const.RIGHT and not self.has_right_wall:
            return True
        if unit_vector == const.UP and not self.has_top_wall:
            return True
        if unit_vector == const.DOWN and not self.has_bottom_wall:
            return True
        return False

    def __repr__(self):
        return f"Cell({self._top_left}, {self._top_right}, {self._bot_left}, {self._bot_right})"
