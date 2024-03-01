from cell import Cell
import constants as const
import time
import random


class Maze:
    def __init__(
        self,
        x1,
        y1,
        num_rows,
        num_cols,
        size_x,
        size_y,
        win=None,
        seed=None,
    ):
        self.x1 = x1
        self.y1 = y1
        self.num_rows = num_rows
        self.num_cols = num_cols
        self.size_x = size_x
        self.size_y = size_y
        self.win = win
        self.seed = seed

        self._cells = list()
        self._create_cells()

    def _create_cells(self):
        for j in range(self.num_cols):
            row = list()
            for i in range(self.num_rows):
                start_x = self.x1 + i * self.size_x
                start_y = self.y1 + j * self.size_y
                row.append(Cell(
                    start_x,
                    start_y,
                    self.size_x,
                    self.size_y,
                    self.win
                ))
            self._cells.append(row)

        if self.win:
            for row in self._cells:
                for cell in row:
                    self._draw_cell(cell)

            self._break_entrance_and_exit()
            self._break_walls_r(0, 0)
            self._reset_cells_visited()

    def _draw_cell(self, cell):
        cell.draw()
        self._animate(0)

    def _animate(self, sleep_val):
        self.win.redraw()
        time.sleep(sleep_val)

    def _break_entrance_and_exit(self):
        top = self._cells[0][0]
        bottom = self._cells[-1][-1]

        top.knock_down_wall(const.UP)
        self._draw_cell(top)
        bottom.knock_down_wall(const.DOWN)
        self._draw_cell(bottom)

    def _is_in_bounds(self, j, i):
        j_in_bounds = j >= 0 and j < self.num_cols
        i_in_bounds = i >= 0 and i < self.num_rows
        return j_in_bounds and i_in_bounds

    def _break_walls_r(self, j, i):
        if self.seed:
            random.seed(self.seed)
        current = self._cells[j][i]
        current.has_visited()
        while True:
            dir_to_visit = list()
            for v in const.UNIT_VECTORS:
                j_, i_ = const.add_uv(j, i, v)
                if not self._is_in_bounds(j_, i_):
                    continue
                if self._cells[j_][i_].visited is False:
                    dir_to_visit.append(v)
            if not dir_to_visit:
                self._draw_cell(current)
                return
            dir = random.choice(dir_to_visit)
            j_, i_ = const.add_uv(j, i, dir)
            destination = self._cells[j_][i_]
            current.knock_down_wall(dir)
            destination.knock_down_wall(const.invert_uv(dir))
            self._draw_cell(current)
            self._draw_cell(destination)
            self._break_walls_r(j_, i_)

    def _reset_cells_visited(self):
        for row in self._cells:
            for cell in row:
                cell.visited = False

    def solve(self):
        return self._solve_r(0, 0)

    def _solve_r(self, j, i):
        self._animate(.005)
        current = self._cells[j][i]
        current.has_visited()
        if current == self._cells[-1][-1]:
            return True
        for v in const.UNIT_VECTORS:
            if not current.is_dir_clear(v):
                continue
            j_, i_ = const.add_uv(j, i, v)
            if not self._is_in_bounds(j_, i_):
                continue
            peeked = self._cells[j_][i_]
            if peeked.visited is True:
                continue
            current.draw_move(peeked)
            if self._solve_r(j_, i_):
                return True
            else:
                current.draw_move(peeked, True)
        return False
