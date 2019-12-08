use std::fs;
use std::io::BufRead;

fn main() {
    let r = fs::File::open("2019_11_allison.txt").expect("could not open file");
    let buf = std::io::BufReader::new(r);

    for line in buf.lines() {
        println!("{}", line.expect("error reading line"));
    }
}
