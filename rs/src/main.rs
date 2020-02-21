extern crate chrono;
extern crate chrono_locale;

use crate::chrono_locale::LocaleDate;
use std::error;
use std::fmt;
use std::fs;
use std::io::BufRead;
use std::str::FromStr;

#[derive(PartialEq, Debug)]
struct DayEntry {
    day: chrono::NaiveDate,
    hours: f64,
    fee: f64,
    meal: f64,
}

#[derive(Debug)]
struct ParseDayEntryError {
    cause: String,
}

impl fmt::Display for ParseDayEntryError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.cause)
    }
}

impl error::Error for ParseDayEntryError {
    fn source(&self) -> Option<&(dyn error::Error + 'static)> {
        Some(self)
    }
}

impl FromStr for DayEntry {
    type Err = Box<dyn error::Error>;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let parts: Vec<&str> = s.split_whitespace().collect();
        if parts.len() != 5 {
            return Err(Box::new(ParseDayEntryError {
                cause: "invalid fields count".to_owned(),
            }));
        }
        Ok(DayEntry {
            day: chrono::NaiveDate::parse_from_str(parts[1], "%d/%m/%Y")
                .or_else(|err| Err(format!("invalid date: {}", err)))?,
            hours: parts[2].parse::<f64>()?,
            fee: parts[3].parse::<f64>()?,
            meal: parts[4].parse::<f64>()?,
        })
    }
}

fn capitalize(s: &str) -> String {
    let mut res = String::with_capacity(s.len());
    for (i, c) in s.chars().enumerate() {
        if i == 0 {
            res.push_str(&c.to_uppercase().to_string());
        } else {
            res.push_str(&c.to_lowercase().to_string());
        }
    }
    res
}

impl ToString for DayEntry {
    fn to_string(&self) -> String {
        let name = capitalize(&self.day.formatl("%a", "fr_FR").to_string());
        return format!(
            "{name} {day} {hours:.02} {fee:.02} {meal:.02}",
            name = name,
            day = self.day.format("%d/%m/%Y"),
            hours = self.hours,
            fee = self.fee,
            meal = self.meal
        );
    }
}

fn read_sheet<I>(lines: I) -> Vec<DayEntry>
where
    I: IntoIterator<Item = String>,
{
    let mut entries = Vec::with_capacity(31);
    for line in lines {
        if line.starts_with('#') || line.is_empty() {
            continue;
        }
        entries.push(line.parse::<DayEntry>().expect(&line));
    }
    entries
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parse_day_entry() {
        let input = "Mer 13/11/2019 4.50 4.10 3.00";
        let expected = DayEntry {
            day: chrono::NaiveDate::from_ymd(2019, 11, 13),
            hours: 4.5,
            fee: 4.1,
            meal: 3.0,
        };
        let actual: DayEntry = input.parse().expect("could not parse day entry");
        assert_eq!(actual, expected);

        let output = actual.to_string();
        assert_eq!(output, input);
    }
}

fn main() {
    let r = fs::File::open("2019_11_allison.txt").expect("could not open file");
    let buf = std::io::BufReader::new(r);
    let entries = read_sheet(buf.lines().map(|l| l.unwrap()));
    for entry in entries {
        println!("{}", entry.to_string());
    }
}
