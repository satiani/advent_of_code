use std::collections::HashMap;
use std::fs;

fn main() {
    let contents = fs::read_to_string("data/input.txt").expect("Error opening file");
    let lines = contents.lines().collect();
    print_checksum(&lines);
    find_similar_ids(&lines);
}

fn print_checksum(lines: &Vec<&str>) {
    let mut two_letter_count = 0;
    let mut three_letter_count = 0;
    for line in lines {
        let mut character_count = HashMap::new();

        for ch in line.chars() {
            let counter = character_count.entry(ch).or_insert(0);
            *counter += 1;
        }

        if character_count.values().any(|&count| count == 2) {
            two_letter_count += 1;
        }

        if character_count.values().any(|&count| count == 3) {
            three_letter_count += 1;
        }
    }

    println!("Two letters: {}", two_letter_count);
    println!("Three letters: {}", three_letter_count);
    println!("Checksum: {}", two_letter_count * three_letter_count);
}

fn find_similar_ids(lines: &Vec<&str>) {
    let list_size = lines.len();
    for i in 0..list_size {
        let base_string = lines[i];
        for j in i + 1..list_size {
            let compared_string = lines[j];
            let mut differences = 0;
            let mut similar_characters = String::with_capacity(20);

            'char: for (c1, c2) in base_string.chars().zip(compared_string.chars()) {
                if c1 != c2 {
                    differences += 1;
                    if differences > 1 {
                        break 'char;
                    }
                } else {
                    similar_characters.push(c1);
                }
            }

            if differences == 1 {
                println!("Similar characters: {}", similar_characters);
                return;
            }
        }
    }
}
