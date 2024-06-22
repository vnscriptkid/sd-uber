import matplotlib.pyplot as plt
import matplotlib.patches as patches

class QuadTree:
    def __init__(self, boundary, capacity):
        self.boundary = boundary  # boundary is a tuple (x, y, width, height)
        self.capacity = capacity  # max number of points each quadrant can hold
        self.points = []
        self.divided = False

    def insert(self, point, direction=None, level = 0):
        if self.divided:
            # If the quadtree is divided, insert the point in one of the children
            return (self.northeast.insert(point, 'northeast', level + 1) or
                    self.northwest.insert(point, 'northwest', level + 1) or
                    self.southeast.insert(point, 'southeast', level + 1) or
                    self.southwest.insert(point, 'southwest', level + 1))
        
        print(f'{level * 5 * '-'}Inserting point {point} in {direction}...')
        # Check if point is out of bounds
        if not self.in_boundary(point):
            print(f'{level * 5 * '-'}Point {point} is out of bounds of {direction}')
            return False

        # If capacity is not exceeded, add point
        if len(self.points) < self.capacity:
            print(f'{level * 5 * '-'}Point {point} added to {self.boundary}')
            print('-' * 20)
            self.points.append(point)
            return True
        else:
            # Otherwise, subdivide
            print(f'{level * 5 * '-'}Capacity exceeded in {self.boundary}')
            if not self.divided:
                self.subdivide(level)

            # Try to insert the point in one of the children
            return (self.northeast.insert(point, 'northeast', level + 1) or
                    self.northwest.insert(point, 'northwest', level + 1) or
                    self.southeast.insert(point, 'southeast', level + 1) or
                    self.southwest.insert(point, 'southwest', level + 1))

    def in_boundary(self, point):
        x, y = point
        x0, y0, w, h = self.boundary
        return x0 <= x < x0 + w and y0 <= y < y0 + h

    def subdivide(self, level=0):
        x, y, w, h = self.boundary
        self.northeast = QuadTree((x + w/2, y, w/2, h/2), self.capacity)
        print(f'{level * 5 * '-'}Created northeast quadrant: {self.northeast.boundary} with capacity {self.capacity}')
        self.northwest = QuadTree((x, y, w/2, h/2), self.capacity)
        print(f'{level * 5 * '-'}Created northwest quadrant: {self.northwest.boundary} with capacity {self.capacity}')
        self.southeast = QuadTree((x + w/2, y + h/2, w/2, h/2), self.capacity)
        print(f'{level * 5 * '-'}Created southeast quadrant: {self.southeast.boundary} with capacity {self.capacity}')
        self.southwest = QuadTree((x, y + h/2, w/2, h/2), self.capacity)
        print(f'{level * 5 * '-'}Created southwest quadrant: {self.southwest.boundary} with capacity {self.capacity}')
        self.divided = True

        # Redistribute existing points to the new quadrants
        points_to_redistribute = self.points.copy()  # Copy the list to avoid modifying it during iteration
        self.points = []  # Clear the current points list
        divisions = [self.northeast, self.northwest, self.southeast, self.southwest]
        for point in points_to_redistribute:
            for division in divisions:
                if division.in_boundary(point):
                    division.insert(point, None, level + 1)
                    break

    def show(self, ax):
        x, y, w, h = self.boundary
        rect = patches.Rectangle((x, y), w, h, linewidth=1, edgecolor='r', facecolor='none')
        ax.add_patch(rect)
        if self.divided:
            self.northeast.show(ax)
            self.northwest.show(ax)
            self.southeast.show(ax)
            self.southwest.show(ax)
        for point in self.points:
            ax.plot(*point, 'bo')

def main():
    boundary = (0, 0, 100, 100)  # x, y, width, height
    qt = QuadTree(boundary, 4)  # boundary and capacity
    points = [(10, 10), (20, 20), (70, 70), (80, 80), (30, 30), (20, 40), (40, 40), (49, 49), (40, 26), (45, 39)]

    for point in points:
        qt.insert(point)

    fig, ax = plt.subplots()
    qt.show(ax)
    plt.xlim(0, 100)
    plt.ylim(0, 100)
    plt.show()

if __name__ == '__main__':
    main()
