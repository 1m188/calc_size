use std::{env, fs, path::Path};

fn calc_size(path: &Path) -> std::io::Result<u64> {
    let metadata = fs::metadata(path)?;

    if metadata.is_file() {
        // 如果是文件，直接返回其大小
        Ok(metadata.len())
    } else if metadata.is_dir() {
        // 如果是目录，递归计算所有内容的大小
        let mut total_size = 0u64;

        // 遍历目录中的所有条目
        for entry in fs::read_dir(path)? {
            let entry = entry?;
            let entry_path = entry.path();

            // 递归计算每个条目的大小并累加
            total_size += calc_size(&entry_path)?;
        }

        Ok(total_size)
    } else {
        // 既不是文件也不是目录（如符号链接等）
        Ok(0)
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();

    // 输出帮助信息
    if args.len() == 2
        && (args[1] == "--help" || args[1] == "--h" || args[1] == "-help" || args[1] == "-h")
    {
        println!("计算文件/文件夹大小\n");
        println!("示例: calc_size_rust [file1] [file2] [folder1] ...\n");
        println!("输出 -1 表示路径错误或其他错误");
        return;
    }

    // 文件路径
    let mut file_paths: Vec<String> = vec![String::new(); args.len() - 1];
    for (i, arg) in args.iter().enumerate() {
        if i == 0 {
            continue;
        }
        file_paths[i - 1] = String::from(arg);
    }

    // 文件大小
    let mut file_sizes: Vec<f64> = vec![-1.; file_paths.len()];

    // 大小单位
    let mut units: Vec<String> = vec![String::from(" B"); file_sizes.len()];

    // 计算大小
    for (i, file_path) in file_paths.iter().enumerate() {
        let path = Path::new(file_path);
        file_sizes[i] = match calc_size(path) {
            Ok(size) => size as f64,
            Err(_) => -1.,
        }
    }

    // 输出格式化结果

    // 文件名对齐
    let mut max_file_name_len = 0;
    for file_path in file_paths.iter() {
        if file_path.len() > max_file_name_len {
            max_file_name_len = file_path.len();
        }
    }
    for i in 0..file_paths.len() {
        if file_paths[i].len() < max_file_name_len {
            let x = max_file_name_len - file_paths[i].len();
            file_paths[i].push_str(&" ".repeat(x));
        }
    }

    // 单位转换
    for i in 0..file_sizes.len() {
        if file_sizes[i] >= 1024. {
            file_sizes[i] /= 1024.;
            units[i] = String::from("KB");

            if file_sizes[i] >= 1024. {
                file_sizes[i] /= 1024.;
                units[i] = String::from("MB");

                if file_sizes[i] >= 1024. {
                    file_sizes[i] /= 1024.;
                    units[i] = String::from("GB");
                }
            }
        }

        if file_sizes[i] < 0. {
            units[i] = String::from("  ");
        }
    }

    // 输出最终结果
    for (i, file_path) in file_paths.iter().enumerate() {
        println!(
            "{} : {:>w$.2} {}",
            file_path,
            file_sizes[i],
            units[i],
            w = 8, // 根据需要调整宽度
        );
    }
}
