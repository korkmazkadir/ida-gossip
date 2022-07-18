#!/usr/bin/env Rscript

source("stats-aggregator.R")

#path <- "/Users/kadir/Desktop/ida-gossip/scripts/stats"
path <- "/Users/kadir/Downloads/experiment-1"

# lists zip files under the folder
zip_files <- list.files(path = path, pattern = "*.zip", full.names = TRUE)

# unzip all files under /temp_stats_314 folder
temp_folder <- paste(path, "/temp_stats_314", sep = "")
for (zip_file in zip_files) {
  unzip(zipfile = zip_file, exdir = temp_folder)
}


# lists directories
directories <- list.dirs(paste(temp_folder, "/stats", sep = ""), full.names = TRUE, recursive = TRUE)

# calculate data frames
calculate_datasets(directories)

# remove temp folder
unlink(temp_folder, recursive = TRUE)

