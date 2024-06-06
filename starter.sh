nohup ./llamafile.exe -m nomic-embed-text-v1.5.Q8_0.gguf -c 8191 --server --nobrowser --embedding --port 9000 > embedder.log &
nohup ./indexer/coffeemaker > indexer.log &