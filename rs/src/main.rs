extern crate chrono;

use std::str::FromStr;
use std::fs;
use std::io::BufRead;
use std::error;
use std::fmt;

#[derive(PartialEq, Debug)]
struct DayEntry {
    day: chrono::NaiveDate,
    hours: f64,
    fee: f64,
    meal: f64,
}

#[derive(Debug)]
struct ParseDayEntryError {
    cause: String
}

impl fmt::Display for ParseDayEntryError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.cause)
    }
}

impl error::Error for ParseDayEntryError {
    fn source(&self) ->  Option<&(dyn error::Error + 'static)> {
        Some(self)
    }
}

impl FromStr for DayEntry {
    type Err = Box<dyn error::Error>;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let parts : Vec<&str> = s.split_whitespace().collect();
        if parts.len() != 5 {
            return Err(Box::new(ParseDayEntryError{cause: "invalid fields count".to_owned()}));
        }
        Ok(DayEntry{
            day: chrono::NaiveDate::parse_from_str(parts[1], "%d/%m/%Y").or_else(|err| {
                Err(format!("invalid date: {}", err))
            })?,
            hours: parts[2].parse::<f64>()?,
            fee: parts[3].parse::<f64>()?,
            meal: parts[4].parse::<f64>()?,
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parse_day_entry() {
        let input = "Mer 13/11/2019 4.50 4.10 3.00";
        let expected = DayEntry{
            day: chrono::NaiveDate::from_ymd(2019, 11, 13),
            hours: 4.5,
            fee: 4.1,
            meal: 3.0,
        };
        let actual : DayEntry = input.parse().expect("could not parse day entry");
        assert_eq!(actual, expected);
    }
}

fn main() {
    let r = fs::File::open("2019_11_allison.txt").expect("could not open file");
    let buf = std::io::BufReader::new(r);

    for line in buf.lines() {
        println!("{}", line.expect("error reading line"));
    }
}
