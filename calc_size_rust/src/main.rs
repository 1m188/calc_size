use astro_float::{BigFloat, Consts, RoundingMode, ctx::Context};
use std::{env, fs, path::Path, sync::LazyLock};

// 大数上下文
static CTX: LazyLock<Context> = LazyLock::new(|| {
    Context::new(
        512,                    // 精度
        RoundingMode::ToEven,   // 舍入
        Consts::new().unwrap(), // 常数缓存
        -100_000,               // 指数下界
        100_000,                // 指数上界
    )
});

// 计算文件/文件夹大小
fn calc_size(path: &Path) -> std::io::Result<BigFloat> {
    let metadata = fs::metadata(path)?;

    if metadata.is_file() {
        // 如果是文件，直接返回其大小
        Ok(BigFloat::from_u64(metadata.len(), CTX.precision()))
    } else if metadata.is_dir() {
        // 如果是目录，递归计算所有内容的大小
        let mut total_size = BigFloat::from_u64(0, CTX.precision());

        // 遍历目录中的所有条目
        for entry in fs::read_dir(path)? {
            let entry = entry?;
            let entry_path = entry.path();

            // 递归计算每个条目的大小并累加
            let x = calc_size(&entry_path)?;
            total_size = total_size.add(&x, CTX.precision(), CTX.rounding_mode());
        }

        Ok(total_size)
    } else {
        // 既不是文件也不是目录（如符号链接等）
        Ok(BigFloat::from_u64(0, CTX.precision()))
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
    let mut file_sizes: Vec<BigFloat> =
        vec![BigFloat::from_f64(-1., CTX.precision()); file_paths.len()];

    // 大小单位
    let mut units: Vec<String> = vec![String::from(" B"); file_sizes.len()];

    // 计算大小
    for (i, file_path) in file_paths.iter().enumerate() {
        let path = Path::new(file_path);
        file_sizes[i] = match calc_size(path) {
            Ok(size) => size,
            Err(_) => BigFloat::from_f64(-1., CTX.precision()),
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
        let two = BigFloat::from_f64(1024., CTX.precision());
        if file_sizes[i] >= two {
            file_sizes[i] = file_sizes[i].div(&two, CTX.precision(), CTX.rounding_mode());
            units[i] = String::from("KB");

            if file_sizes[i] >= two {
                file_sizes[i] = file_sizes[i].div(&two, CTX.precision(), CTX.rounding_mode());
                units[i] = String::from("MB");

                if file_sizes[i] >= two {
                    file_sizes[i] = file_sizes[i].div(&two, CTX.precision(), CTX.rounding_mode());
                    units[i] = String::from("GB");
                }
            }
        }

        if file_sizes[i] < BigFloat::from_f64(0., CTX.precision()) {
            units[i] = String::from("  ");
        }
    }

    // 输出最终结果
    for (i, file_path) in file_paths.iter().enumerate() {
        // 手动转换file_size
        let s = file_sizes[i].to_string();
        let li: Vec<&str> = s.split('e').collect();
        let e = li[1].parse::<i32>().unwrap();
        let x = li[0].parse::<f64>().unwrap() * 10_f64.powi(e);

        println!("{} : {:>8.2} {}", file_path, x, units[i]);
    }
}
