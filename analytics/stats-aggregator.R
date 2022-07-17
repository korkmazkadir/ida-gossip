#!/usr/bin/env Rscript

library("rjson")
library(dplyr)

standart_error <-function(df){
  return(sd(df$Value) / sqrt(nrow(df)))
}

confidence_interval <- function(df) {
  z_val <- qnorm(.025, lower.tail = FALSE)
  ci_val <- z_val * standart_error(df)
  return(ci_val)
}


create_first_chunk_df <- function(config, stats_df) {
  result_df <- as.data.frame(config)
  df <- stats_df %>% filter(Event == "FIRST_CHUNK_RECEIVED")
  stats <- boxplot.stats(df$Value)$stats

  mean <- mean(df$Value)
  ci <- confidence_interval(df)
  
  result_df["Min"] <- c(stats[1])
  result_df["FirstQuartile"] <- c(stats[2])
  result_df["Median"] <- c(stats[3])
  result_df["ThirdQuartile"] <- c(stats[4])
  result_df["Max"] <- c(stats[5])
  result_df["RowCount"] <- c( nrow(df) )
  result_df["MeanLowerBound"] <- mean - ci
  result_df["Mean"] <- mean
  result_df["MeanUpperBound"] <- mean + ci
  
  return(result_df)
}


create_message_received_df <- function(config, stats_df) {
  result_df <- as.data.frame(config)
  df <- stats_df %>% filter(Event == "MESSAGE_RECEIVED")
  stats <- boxplot.stats(df$Value)$stats
  
  mean <- mean(df$Value)
  ci <- confidence_interval(df)
  
  result_df["Min"] <- c(stats[1])
  result_df["FirstQuartile"] <- c(stats[2])
  result_df["Median"] <- c(stats[3])
  result_df["ThirdQuartile"] <- c(stats[4])
  result_df["Max"] <- c(stats[5])
  result_df["RowCount"] <- c( nrow(df) )
  result_df["MeanLowerBound"] <- mean - ci
  result_df["Mean"] <- mean
  result_df["MeanUpperBound"] <- mean + ci
  
  return(result_df)
}


create_network_usage_df <- function(config, stats_df) {
  result_df <- as.data.frame(config)
  df <- stats_df %>% filter(Event == "NETWORK_USAGE")
  
  # TODO: is it a good way to solve the problem
  # int overflow problem 
  df$Value <- as.numeric(df$Value)
  
  typeof(df$Value)
  
  stats <- boxplot.stats(df$Value)$stats
  
  mean <- mean(df$Value)
  ci <- confidence_interval(df)
  
  result_df["Min"] <- c(stats[1])
  result_df["FirstQuartile"] <- c(stats[2])
  result_df["Median"] <- c(stats[3])
  result_df["ThirdQuartile"] <- c(stats[4])
  result_df["Max"] <- c(stats[5])
  result_df["RowCount"] <- c( nrow(df) )
  result_df["MeanLowerBound"] <- mean - ci
  result_df["Mean"] <- mean
  result_df["MeanUpperBound"] <- mean + ci
  
  return(result_df)
}


# iterates over directories
calculate_datasets <- function(directories) {
  print("Processing...")

  first_chunk_delivery_df <- data.frame(matrix(ncol = 13, nrow = 0))
  message_received_df <- data.frame(matrix(ncol = 11, nrow = 0))
  network_usage_df <- data.frame(matrix(ncol = 11, nrow = 0))
  
  for (directory in directories) {
    config_file <- paste(directory, "/config.json", sep = "")
    if (file.exists(config_file) == FALSE) {
      next
    }

    print(directory)

    stats_file <- paste(directory, "/stats.log", sep = "")

    config <- fromJSON(file = config_file)
    stats <- read.csv(stats_file, sep = "\t", header = FALSE)
    colnames(stats) <- c("NodeID", "Round", "Event", "Value")

    # first chunk delivery
    df <- create_first_chunk_df(config, stats)
    first_chunk_delivery_df <- rbind(first_chunk_delivery_df, df)

    # message delivery
    df <- create_message_received_df(config, stats)
    message_received_df <- rbind(message_received_df, df)

    # queue size
    df <- create_network_usage_df(config, stats)
    network_usage_df <- rbind(network_usage_df, df)
  }

  # write data frames
  write.table(first_chunk_delivery_df, file = paste(path, "/first_chunk_delivery_df.tsv", sep = ""), quote = FALSE, sep = "\t", col.names = TRUE, row.names = FALSE)
  write.table(message_received_df, file = paste(path, "/message_received_df.tsv", sep = ""), quote = FALSE, sep = "\t", col.names = TRUE, row.names = FALSE)
  write.table(network_usage_df, file = paste(path, "/network_usage_df.tsv", sep = ""), quote = FALSE, sep = "\t", col.names = TRUE, row.names = FALSE)
}

#path <- "/home/kadir/Desktop/ida_connection_count_effect"

# lists directories
#directories <- list.dirs(path = path, full.names = TRUE, recursive = TRUE)

# calculate data frames
#calculate_datasets(directories)
