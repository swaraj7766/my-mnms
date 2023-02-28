# mininet

path: root of mnms

### Manual

1. ### install

   ```sh
   sudo apt-get install mininet
   sudo apt install net-tools
   sudo apt-get install xterm
   ```

2. ### run minnet and two of host 

   ```sh
   sudo mn
   ```

   

3. ### open node of  h1 and h2 

   ```sh
   xterm h1 h2
   ```

   

4. ### run simulator at h2

   ```sh
   ./pkg/simulator/bin/simulator.exe run -d
   ```

   

5. ### run mnms service at h1

   ```sh
   ./mnmsctl -n myname2 -s
   ```

   

6. ### exit min

   ```sh
   exit
   ```

   

