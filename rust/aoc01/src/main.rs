use std::collections::HashSet;
use std::fs;
use std::io::Result;

fn main() -> Result<()> {
    let path = "data/input.txt";
    let lines = fs::read_to_string(path)?;
    let values = lines.lines().map(|v| v.parse::<i32>().unwrap());
    println!("Sum: {}", values.clone().sum::<i32>());

    let mut seen = HashSet::<i32>::default();
    let mut total: i32 = 0;
    seen.insert(0);

    for value in values.clone().cycle() {
        total += value;
        if !seen.insert(total) {
            println!("Repeated Frequency: {}", total);
            return Ok(());
        }
    }

    Ok(())
}
