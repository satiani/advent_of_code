use std::fs;
use std::io::Result;
use std::collections::HashSet;

fn main() -> Result<()> {
    let path = "data/input.txt";
    let lines = fs::read_to_string(path)?;
    let mut total: i32 = 0;
    for line in lines.lines() {
        let change: i32 = line.parse().unwrap();
        total += change
    }
    println!("Sum: {}", total);

    let mut seen = HashSet::new();
    total = 0;
    seen.insert(0);
    loop {
        for line in lines.lines() {
            let change: i32 = line.parse().unwrap();
            total += change;
            if seen.contains(&total) {
                println!("Repeated Frequency: {}", total);
                return Ok(());
            }
            seen.insert(total);
        }
    }
}
