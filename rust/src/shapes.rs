use std::collections::{HashSet};
use image::{ImageBuffer, Rgb};
use crate::{MULTIPLIER, OFFSET};

#[derive(PartialEq, Eq, Hash, Clone)]
pub enum Direction {
    Up,
    Down,
    Left,
    Right,
}

#[derive(Eq, Hash, PartialEq, Clone, Copy)]
pub struct Point (pub i32, pub i32);

#[derive(PartialEq, Eq, Hash)]
pub struct Line {
    pub start: Point,
    pub end: Point,
    pub shown: bool,
    pub color: Rgb<u8>,
}

#[derive(PartialEq, Eq, Hash)]
pub struct Field<'a> {
    pub borders: Vec<&'a mut Line>,
    pub point: Point,
}

pub struct Labyrinth<'a> {
    pub fields: Vec<Vec<&'a mut Field<'a>>>,
    pub size: usize,
    lines: HashSet<&'a mut Line>,
}

impl Line {
    pub fn new(point: Point, direction: Direction, color: Option<Rgb<u8>>) -> Line {
        let color = color.unwrap_or(Rgb([0u8, 0u8, 0u8]));
        let (dx, dy) = direction.relative();

        let start_x = (point.0 as isize + dx) as usize * MULTIPLIER + OFFSET;
        let start_y = (point.1 as isize + dy) as usize * MULTIPLIER + OFFSET;
        let end_x = (point.0 as isize + dx) as usize * MULTIPLIER + OFFSET;
        let end_y = (point.1 as isize + dy) as usize * MULTIPLIER + OFFSET;

        Line {
            start: Point::new(start_x, start_y),
            end: Point::new(end_x, end_y),
            shown: true,
            color,
        }
    }

    pub fn hide(&mut self) {
        self.shown = false;
    }

    pub fn draw(&self, img: &mut ImageBuffer<Rgb<u8>, Vec<u8>>) {
        if self.shown {
            imageproc::drawing::draw_line_segment_mut(img, self.start.tuple(), self.end.tuple(), self.color);
        }
    }
}

impl Point {
    pub fn new(x: usize, y: usize) -> Point {
        Point(x as i32, y as i32)
    }

    pub fn tuple(&self) -> (f32, f32) {
        (self.0 as f32, self.1 as f32)
    }
}

impl<'a> Field<'a> {
    pub fn new(point: Point, borders: Vec<&mut Line>) -> Field {
        Field {
            borders,
            point
        }
    }

    pub fn move_to(&'a self, labyrinth: &'a Labyrinth<'a>, direction: &Direction) -> Option<&mut Field> {
        let (dx, dy) = direction.relative();
        let x = self.point.0 as isize + dx;
        let y = self.point.1 as isize + dy;

        if 0 <= x && x < labyrinth.size as isize && 0 <= y && y < labyrinth.size as isize {
            Some(labyrinth.fields[y as usize][x as usize])
        } else {
            None
        }
    }
}

impl Labyrinth<'_> {
    pub fn new<'a>(size: usize) -> Labyrinth<'a> {
        let mut fields: Vec<Vec<&mut Field>> = Vec::with_capacity(size);
        let mut lines: HashSet<&mut Line> = HashSet::new();

        for y in 0..size {
            let mut row: Vec<&mut Field> = Vec::with_capacity(size);
            for x in 0..size {
                let point = Point::new(x, y);
                let mut l: Vec<Line> = [Direction::Up, Direction::Down, Direction::Left, Direction::Right]
                    .iter()
                    .map(|d| Line::new(point, d.clone(), None))
                    .collect();
                
                let mut l: Vec<& mut Line> = l.iter_mut().collect();

                if y != 0 {
                    l[Direction::Up.idx()] = fields[y-1][x].borders[Direction::Down.idx()];
                }
                if x != 0 {
                    l[Direction::Left.idx()] = fields[y][x-1].borders[Direction::Right.idx()];
                }

                row.push(&mut Field::new(point, l));
                for line in l {
                    lines.insert(line);
                }
            }
            fields.push(row);
        }
        
        Labyrinth {
            fields,
            size,
            lines,
        }
    }

    pub fn draw(&self, img: &mut ImageBuffer<Rgb<u8>, Vec<u8>>) {
        for line in self.lines.iter() {
            line.draw(img);
        }
    }
}

impl Direction {
    pub fn relative(&self) -> (isize, isize) {
        match self {
            Direction::Up => (0, -1),
            Direction::Down => (0, 1),
            Direction::Left => (-1, 0),
            Direction::Right => (1, 0),
        }
    }

    pub fn iter() -> std::slice::Iter<'static, Direction> {
        static DIRECTIONS: [Direction; 4] = [Direction::Up, Direction::Down, Direction::Left, Direction::Right];
        DIRECTIONS.iter()
    }

    pub fn idx(&self) -> usize {
        match self {
            Direction::Up => 0,
            Direction::Down => 1,
            Direction::Left => 2,
            Direction::Right => 3,
        }
    }
}