extern crate chrono;
#[macro_use]
extern crate serde_scan;

use chrono::{NaiveDateTime, Timelike};
use std::collections::HashMap;
use std::fs;
use std::str::Lines;

#[derive(Debug)]
enum EntryType {
    StartedShift(u16),
    StartedSleeping,
    WokeUp,
}

#[derive(Debug)]
struct LogEntry {
    entry_type: EntryType,
    date_time: NaiveDateTime,
}

struct GuardInfo {
    total_sleep_duration: u16,
    minute_sleep_frequency: [u16; 60],
}

type GuardStats = HashMap<u16, GuardInfo>;


fn main() {
    let contents = fs::read_to_string("data/input.txt").unwrap();
    if let Some(log_entries) = gather_log_entries(contents.lines()) {
        let guard_stats = generate_guard_stats(log_entries);
        print_sleepiest_guard(&guard_stats);
        print_most_frequently_slept_minute(&guard_stats);
    }
}

fn gather_log_entries(lines: Lines) -> Option<Vec<LogEntry>> {
    let mut log_entries: Vec<LogEntry> = Vec::new();
    for line in lines {
        let strings: Vec<&str> = scan!("[{}] {}" <- line).unwrap();
        if let [date, time, entry_start] = strings[0..3] {
            let datetime =
                NaiveDateTime::parse_from_str(&format!("{} {}", date, time), "%Y-%m-%d %H:%M")
                    .unwrap();

            let entry_type: EntryType = match entry_start {
                "wakes" => EntryType::WokeUp,
                "falls" => EntryType::StartedSleeping,
                "Guard" => {
                    let guard_id = strings[3].trim_matches('#').parse::<u16>().unwrap();
                    EntryType::StartedShift(guard_id)
                }
                _ => unreachable!(),
            };

            log_entries.push(LogEntry {
                entry_type: entry_type,
                date_time: datetime,
            });
        }
    }

    if log_entries.len() > 0 {
        log_entries.sort_by_key(|entry: &LogEntry| entry.date_time);
        return Some(log_entries);
    }
    None
}

fn generate_guard_stats(log_entries: Vec<LogEntry>) -> GuardStats {
    // map between guard id to tuple of total minutes slept and frequency
    // sleep times during midnight minutes
    let mut guard_stats: GuardStats = HashMap::default();
    let mut current_guard_id: u16 = 0;
    let mut current_sleep_start_time: NaiveDateTime = NaiveDateTime::from_timestamp(0, 0);
    for entry in log_entries {
        match entry.entry_type {
            EntryType::StartedShift(guard_id) => current_guard_id = guard_id,
            EntryType::StartedSleeping => current_sleep_start_time = entry.date_time,
            EntryType::WokeUp => {
                let guard_info = guard_stats.entry(current_guard_id).or_insert(GuardInfo {
                    total_sleep_duration: 0,
                    minute_sleep_frequency: [0u16; 60],
                });
                let sleep_start_min = current_sleep_start_time.minute() as u16;
                let sleep_end_min = entry.date_time.minute() as u16;
                guard_info.total_sleep_duration += sleep_end_min - sleep_start_min;
                for i in sleep_start_min..sleep_end_min {
                    guard_info.minute_sleep_frequency[i as usize] += 1;
                }
            }
        }
    }
    guard_stats
}

fn print_sleepiest_guard(guard_stats: &GuardStats) {
    let mut longest_sleep_duration = 0u16;

    let mut longest_sleep_guard_id = 0u16;
    for (guard_id, guard_info) in guard_stats {
        if guard_info.total_sleep_duration > longest_sleep_duration {
            longest_sleep_duration = guard_info.total_sleep_duration;
            longest_sleep_guard_id = *guard_id;
        }
    }

    let sleepiest_guard_info = guard_stats.get(&longest_sleep_guard_id).unwrap();
    let minute_frequency: [u16; 60] = sleepiest_guard_info.minute_sleep_frequency;
    let mut max_frequency: u16 = 0;
    let mut max_frequency_index: usize = 0;
    for (i, frequency) in minute_frequency.iter().enumerate() {
        if *frequency > max_frequency {
            max_frequency = *frequency;
            max_frequency_index = i;
        }
    }
    println!("Sleepiest guard id is: {}", longest_sleep_guard_id);
    println!(
        "Sleepiest minute for that guard is: {}",
        max_frequency_index
    );
}

fn print_most_frequently_slept_minute(guard_stats: &GuardStats) {
    let mut most_frequently_slept_minute = 0u16;
    let mut highest_frequency = 0u16;
    let mut guard_with_most_frequently_slept_minute = 0u16;

    for (guard_id, guard_info) in guard_stats {
        for (i, frequency) in guard_info.minute_sleep_frequency.iter().enumerate() {
            if *frequency > highest_frequency{
                highest_frequency = *frequency;
                most_frequently_slept_minute = i as u16;
                guard_with_most_frequently_slept_minute = *guard_id;
            }
        }
    }

    println!(
        "Sleepiest guard id: {}",
        guard_with_most_frequently_slept_minute
    );
    println!(
        "Highest frequency: {}",
        highest_frequency
    );
    println!(
        "Most frequently slept minute: {}",
        most_frequently_slept_minute
    );
}
