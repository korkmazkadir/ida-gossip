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

create_send_time_df <- function(config, stats_df) {
  result_df <- as.data.frame(config)
  df <- stats_df %>% filter(Event == "MEAN_SEND_TIME")
  
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
  send_time_df <- data.frame(matrix(ncol = 11, nrow = 0))

  stats <- data.frame(matrix(ncol = 4, nrow = 0))

  config <- data.frame()

  for (directory in directories) {
    config_file <- paste(directory, "/config.json", sep = "")
    if (file.exists(config_file) == FALSE) {
      next
    }

    print(directory)

    stats_file <- paste(directory, "/stats.log", sep = "")

    config <- fromJSON(file = config_file)
    run_stats <- read.csv(stats_file, sep = "\t", header = FALSE)
    stats <- rbind(stats, run_stats)

  }

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

    # queue size
    df <- create_send_time_df(config, stats)
    send_time_df <- rbind(send_time_df, df)
    
    
  print(file.exists("/first_chunk_delivery_df.tsv"))

  # write data frames

  path_chunk <- paste(path, "/first_chunk_delivery_df.tsv", sep = "")
  path_message <- paste(path, "/message_received_df.tsv", sep = "")
  path_network <- paste(path, "/network_usage_df.tsv", sep = "")
  path_send_time <- paste(path, "/send_time_df.tsv", sep = "")

  write.table(first_chunk_delivery_df, file = path_chunk, quote = FALSE, sep = "\t", col.names = !file.exists(path_chunk), row.names = FALSE, append = TRUE)

  write.table(message_received_df, file = path_message, quote = FALSE, sep = "\t", col.names = !file.exists(path_message), row.names = FALSE, append = TRUE)

  write.table(network_usage_df, file = path_network, quote = FALSE, sep = "\t", col.names = !file.exists(path_network), row.names = FALSE, append = TRUE)

  write.table(send_time_df, file = path_send_time, quote = FALSE, sep = "\t", col.names = !file.exists(path_send_time), row.names = FALSE, append = TRUE)
  
  
}

#path <- "/home/kadir/Desktop/ida_connection_count_effect"

# lists directories
#directories <- list.dirs(path = path, full.names = TRUE, recursive = TRUE)

# calculate data frames
#calculate_datasets(directories)
