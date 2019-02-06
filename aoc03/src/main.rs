#[macro_use]
extern crate serde_scan;
#[macro_use]
extern crate serde_derive;

use std::fs;
use std::str::Lines;

#[derive(Deserialize, Debug, PartialEq, Copy, Clone)]
struct Claim {
    id: u16,
    left: u16,
    top: u16,
    width: u16,
    height: u16,
}

fn main() {
    let contents = fs::read_to_string("data/input.txt").expect("Error reading file");
    match gather_claims(contents.lines()) {
        Some((claims, (width, height))) => get_overlapping_inches_count(&claims, width, height),
        None => {
            println!("Unable to get claims from input");
            return;
        }
    }
}

fn gather_claims(lines: Lines) -> Option<(Vec<Claim>, (u16, u16))> {
    let mut claims: Vec<Claim> = Vec::new();
    let mut max_height: u16 = 0;
    let mut max_width: u16 = 0;
    for line in lines {
        let claim: Claim = scan!("#{} @ {},{}: {}x{}" <- line).expect("Failed!");
        claims.push(claim);

        if claim.height + claim.top > max_height {
            max_height = claim.height + claim.top;
        }
        if claim.width + claim.left > max_width {
            max_width = claim.width + claim.left;
        }
    }

    if claims.len() > 0 {
        return Some((claims, (max_width + 1, max_height + 1)));
    }

    None
}

fn get_overlapping_inches_count(claims: &Vec<Claim>, width: u16, height: u16) {
    let mut square_inches = vec![0u8; width as usize * height as usize];
    let mut inches_with_overlapping_claims = 0u32;
    for claim in claims {
        for i in 0..claim.width {
            for j in 0..claim.height {
                let index: usize = (claim.top + j) as usize * 1000 + (claim.left + i) as usize;
                let new_count = square_inches[index].saturating_add(1);
                if new_count == 2 {
                    inches_with_overlapping_claims += 1;
                }
                square_inches[index] = new_count;
            }
        }
    }

    println!("Overlapping inches: {}", inches_with_overlapping_claims);

    'claim: for claim in claims {
        for i in 0..claim.width {
            for j in 0..claim.height {
                let index: usize = (claim.top + j) as usize * 1000 + (claim.left + i) as usize;
                if square_inches[index] > 1 {
                    continue 'claim;
                }
            }
        }
        println!("This claim does not overlap: {:?}", claim);
    }
}
