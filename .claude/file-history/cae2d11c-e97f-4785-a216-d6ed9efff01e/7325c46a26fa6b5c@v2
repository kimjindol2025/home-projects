// I/O Functions (15)
// File and console operations

use crate::core::Value;
use std::fs::{File, OpenOptions};
use std::io::{self, BufReader, BufRead, Read, Write};
use std::collections::HashMap;

// Global file handles for tracking open files
thread_local! {
    static FILE_HANDLES: std::cell::RefCell<HashMap<String, Box<dyn std::any::Any>>> =
        std::cell::RefCell::new(HashMap::new());
}

/// Print arguments without newline
pub fn print(args: Vec<Value>) -> Result<Value, String> {
    for arg in args {
        print!("{}", arg);
    }
    Ok(Value::Null)
}

/// Print arguments with newline
pub fn println(args: Vec<Value>) -> Result<Value, String> {
    for arg in args {
        println!("{}", arg);
    }
    Ok(Value::Null)
}

/// Read a line from stdin
pub fn input(_args: Vec<Value>) -> Result<Value, String> {
    let mut buffer = String::new();
    io::stdin().read_line(&mut buffer)
        .map_err(|e| format!("Input error: {}", e))?;
    Ok(Value::String(buffer.trim().to_string()))
}

/// Open a file for reading (mode="r") or writing (mode="w")
pub fn open(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("open requires 2 arguments: filename, mode".to_string());
    }

    let filename = args[0].to_string();
    let mode = args[1].to_string();

    match mode.as_str() {
        "r" => {
            match File::open(&filename) {
                Ok(file) => {
                    let reader = Box::new(BufReader::new(file)) as Box<dyn std::any::Any>;
                    FILE_HANDLES.with(|fh| {
                        fh.borrow_mut().insert(filename.clone(), reader);
                    });
                    Ok(Value::String(filename))
                }
                Err(e) => Err(format!("Cannot open file: {}", e))
            }
        }
        "w" => {
            match OpenOptions::new()
                .write(true)
                .create(true)
                .truncate(true)
                .open(&filename)
            {
                Ok(file) => {
                    let writer = Box::new(file) as Box<dyn std::any::Any>;
                    FILE_HANDLES.with(|fh| {
                        fh.borrow_mut().insert(filename.clone(), writer);
                    });
                    Ok(Value::String(filename))
                }
                Err(e) => Err(format!("Cannot create file: {}", e))
            }
        }
        _ => Err(format!("Unknown mode: {}", mode))
    }
}

/// Close a file handle
pub fn close(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("close requires filename argument".to_string());
    }

    let filename = args[0].to_string();
    FILE_HANDLES.with(|fh| {
        fh.borrow_mut().remove(&filename);
    });
    Ok(Value::Null)
}

/// Read entire file content as string
pub fn read_file(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("read_file requires filename argument".to_string());
    }

    let filename = args[0].to_string();
    match std::fs::read_to_string(&filename) {
        Ok(content) => Ok(Value::String(content)),
        Err(e) => Err(format!("Cannot read file: {}", e))
    }
}

/// Write string to file
pub fn write_file(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("write_file requires 2 arguments: filename, content".to_string());
    }

    let filename = args[0].to_string();
    let content = args[1].to_string();

    match std::fs::write(&filename, &content) {
        Ok(_) => Ok(Value::Number(content.len() as f64)),
        Err(e) => Err(format!("Cannot write file: {}", e))
    }
}

/// Append string to file
pub fn append_file(args: Vec<Value>) -> Result<Value, String> {
    if args.len() < 2 {
        return Err("append_file requires 2 arguments: filename, content".to_string());
    }

    let filename = args[0].to_string();
    let content = args[1].to_string();

    match OpenOptions::new()
        .create(true)
        .append(true)
        .open(&filename)
    {
        Ok(mut file) => {
            match file.write_all(content.as_bytes()) {
                Ok(_) => Ok(Value::Number(content.len() as f64)),
                Err(e) => Err(format!("Write error: {}", e))
            }
        }
        Err(e) => Err(format!("Cannot open file: {}", e))
    }
}

/// Check if file exists
pub fn file_exists(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("file_exists requires filename argument".to_string());
    }

    let filename = args[0].to_string();
    Ok(Value::Bool(std::path::Path::new(&filename).exists()))
}

/// Delete a file
pub fn delete_file(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("delete_file requires filename argument".to_string());
    }

    let filename = args[0].to_string();
    match std::fs::remove_file(&filename) {
        Ok(_) => Ok(Value::Null),
        Err(e) => Err(format!("Cannot delete file: {}", e))
    }
}

/// Get file size in bytes
pub fn file_size(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("file_size requires filename argument".to_string());
    }

    let filename = args[0].to_string();
    match std::fs::metadata(&filename) {
        Ok(metadata) => Ok(Value::Number(metadata.len() as f64)),
        Err(e) => Err(format!("Cannot get file size: {}", e))
    }
}

/// List files in directory
pub fn list_dir(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("list_dir requires directory path argument".to_string());
    }

    let path = args[0].to_string();
    match std::fs::read_dir(&path) {
        Ok(entries) => {
            let mut files = Vec::new();
            for entry in entries {
                if let Ok(entry) = entry {
                    if let Some(name) = entry.file_name().to_str() {
                        files.push(Value::String(name.to_string()));
                    }
                }
            }
            Ok(Value::array(files))
        }
        Err(e) => Err(format!("Cannot list directory: {}", e))
    }
}

/// Create a directory
pub fn mkdir(args: Vec<Value>) -> Result<Value, String> {
    if args.is_empty() {
        return Err("mkdir requires directory path argument".to_string());
    }

    let path = args[0].to_string();
    match std::fs::create_dir_all(&path) {
        Ok(_) => Ok(Value::Null),
        Err(e) => Err(format!("Cannot create directory: {}", e))
    }
}

/// Print to stderr
pub fn eprint(args: Vec<Value>) -> Result<Value, String> {
    for arg in args {
        eprint!("{}", arg);
    }
    Ok(Value::Null)
}

/// Print to stderr with newline
pub fn eprintln(args: Vec<Value>) -> Result<Value, String> {
    for arg in args {
        eprintln!("{}", arg);
    }
    Ok(Value::Null)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_print() {
        let result = print(vec![Value::String("hello".to_string())]);
        assert!(result.is_ok());
    }

    #[test]
    fn test_println() {
        let result = println(vec![Value::String("hello".to_string())]);
        assert!(result.is_ok());
    }

    #[test]
    fn test_eprint() {
        let result = eprint(vec![Value::String("error".to_string())]);
        assert!(result.is_ok());
    }

    #[test]
    fn test_eprintln() {
        let result = eprintln(vec![Value::String("error".to_string())]);
        assert!(result.is_ok());
    }

    #[test]
    fn test_write_and_read_file() {
        let filename = "/tmp/test_io.txt";
        let content = "test content";

        // Write file
        let write_result = write_file(vec![
            Value::String(filename.to_string()),
            Value::String(content.to_string()),
        ]);
        assert!(write_result.is_ok());

        // Read file
        let read_result = read_file(vec![Value::String(filename.to_string())]);
        assert!(read_result.is_ok());
        if let Ok(Value::String(text)) = read_result {
            assert_eq!(text, content);
        }

        // Clean up
        let _ = delete_file(vec![Value::String(filename.to_string())]);
    }

    #[test]
    fn test_file_exists() {
        let filename = "/tmp/test_exist.txt";

        // Should not exist initially
        let result = file_exists(vec![Value::String(filename.to_string())]);
        assert!(result.is_ok());
        assert!(!result.unwrap().is_truthy());

        // Create file
        let _ = write_file(vec![
            Value::String(filename.to_string()),
            Value::String("test".to_string()),
        ]);

        // Should exist now
        let result = file_exists(vec![Value::String(filename.to_string())]);
        assert!(result.is_ok());
        assert!(result.unwrap().is_truthy());

        // Clean up
        let _ = delete_file(vec![Value::String(filename.to_string())]);
    }

    #[test]
    fn test_file_size() {
        let filename = "/tmp/test_size.txt";
        let content = "test";

        // Write file
        let _ = write_file(vec![
            Value::String(filename.to_string()),
            Value::String(content.to_string()),
        ]);

        // Get size
        let result = file_size(vec![Value::String(filename.to_string())]);
        assert!(result.is_ok());
        if let Ok(Value::Number(size)) = result {
            assert_eq!(size, content.len() as f64);
        }

        // Clean up
        let _ = delete_file(vec![Value::String(filename.to_string())]);
    }

    #[test]
    fn test_mkdir() {
        let path = "/tmp/test_mkdir";

        // Create directory
        let result = mkdir(vec![Value::String(path.to_string())]);
        assert!(result.is_ok());

        // Check if exists
        let exists = std::path::Path::new(path).exists();
        assert!(exists);

        // Clean up
        let _ = std::fs::remove_dir(path);
    }
}
