mod shapes;

use std::collections::{HashMap, HashSet};
use image::{ImageBuffer, Rgb};
use imageproc::drawing::draw_line_segment_mut;
use crate::shapes::{Direction, Field, Labyrinth};

pub const GRID_SIZE: usize = 10;
pub const MULTIPLIER: usize = 25;
pub const OFFSET: usize = 50;

fn main() {
    let start = std::time::Instant::now();
    let mut img = ImageBuffer::new(100, 100);
    for pixel in img.pixels_mut() {
        *pixel = Rgb([255u8, 255u8, 255u8]);
    }

    let mut labyrinth = Labyrinth::new(10);
    labyrinth.fields[0][0].borders[Direction::Left.idx()].hide();
    labyrinth.fields[labyrinth.size-1][labyrinth.size-1].borders[Direction::Right.idx()].hide();

    let mut been_to: HashSet<&Field> = HashSet::with_capacity(labyrinth.size * labyrinth.size);
    let mut stack: Vec<&mut Field> = Vec::with_capacity(labyrinth.size * labyrinth.size);
    let mut pointer = &labyrinth.fields[0][0];

    while been_to.len() < labyrinth.size * labyrinth.size {
        let mut available: HashMap<Direction, &mut Field> = HashMap::new();
        for direction in Direction::iter() {
            if let Some(potential_pointer) = pointer.move_to(&labyrinth, direction) {
                if !been_to.contains(potential_pointer) {
                    available.insert(direction.clone(), potential_pointer);
                }
            }
        }

        been_to.insert(pointer);
        if available.is_empty() {
            pointer = &stack.pop().unwrap();
            continue
        } else if available.len() > 1 {
            stack.push(*pointer);
        }

        let direction = available.keys().nth(rand::random::<usize>() % available.len()).unwrap();
        pointer.borders[direction.idx()].hide();
        pointer = &available[direction];
    }

    draw_line_segment_mut(&mut img, (50.0, 50.0), (70.0, 70.0), Rgb([0u8, 0u8, 0u8]));

    img.save("labyrinth.png").unwrap();
    println!("Time: {:?}", start.elapsed());
}