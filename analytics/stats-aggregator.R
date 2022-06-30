
library("rjson")

create_first_chunk_df <- function(config, stats_df){
  result_df <- as.data.frame(config)
  df <- stats_df %>% filter( Event == "FIRST_CHUNK_RECEIVED" )
  stats <- boxplot.stats(df$Value)$stats

  result_df["Min"] <- c( stats[1] )
  result_df["FirstQuartile"] <- c( stats[2] )
  result_df["Median"] <- c( stats[3] )
  result_df["ThirdQuartile"] <- c( stats[4]  )
  result_df["Max"] <- c( stats[5]  )
  
  return(result_df)
}


create_message_received_df <- function(config, stats_df){
  result_df <- as.data.frame(config)
  df <- stats_df %>% filter( Event == "MESSAGE_RECEIVED" )
  
  return(result_df)
}


create_queue_size_df <- function(config, stats_df){
  result_df <- as.data.frame(config)
  df <- stats_df %>% filter( Event == "QUEUE_LENGTH" )
  
  return(result_df)
}


# iterates over directories
calculate_datasets <-function(config, directories){
    file.exists("config.json")
    config <- fromJSON(file = "config.json")
    stats <- read.csv("./stats.log", sep = "\t", header = FALSE)
    colnames(stats) <- c('NodeID','Round','Event', 'Value')
}


# lists directories
directories <- list.dirs(path = ".", full.names = TRUE, recursive = TRUE)

print(directories)

# iterates over directories
#calculate_datasets(directories)

