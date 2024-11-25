# WASA TEXT

Spostarsi su una cartella (es. Documents). Eseguire il comando: 
```bash 
git clone https://github.com/sapienzaapps/fantastic-coffee-decaffeinated.git
```

Entrare sempre da terminale nella cartella fantastic-coffee-decaffeinated ed eseguire il comando:
```bash
go mod edit -module <nome_modulo>
```
Dopo aver eseguito il comando entrare nei file .go all'interno delle varie cartelle: 
- ./service/api 
- ./service/database
- ./cmd
e cambiare il modulo principale da git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated al nome scelto nel comando 
precedente. Se si trova nei file il modulo + path (esempio: git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api) bisogna cambiare solo la parte iniziale prima del path (Nell'esempio avrete nome_modulo/service/api).

Dopo aver cambiato tutti i moduli, eseguire il comando:
```bash
go mod vendor 
```
Dopo aver eseguito il comando non deve dare errori/scritte, se da delle scritte non sono state cambiate tutte le occorrenze del modulo nei file delle cartelle specificate, ricontrollare i file e rieseguire il comando vendor. Dopo aver fatto questo passaggio 
modificare il file api.yaml in /doc. 

Se avete già caricato l'api su github, andate sulla repo github, clonate (se non lo avete già fatto) la vostra repo, dopo averla clonata spostate il contenuto della cartella fantastic-coffee-decaffeinated (tranne la cartella doc) nella vostra cartella appena clonata, dopo aver fatto questo passaggio, eseguite il comando:
```bash
git add .
git commit -m "First commit"
git push
```

Se chiede l'accesso, mettere username e come password mettete un token generato da github. Per generare il token andate su github, cliccate sulla vostra immagine in alto a destra, cliccate su settings, andate su Developer settings, cliccate su personal access tokens, token classi, generate token classic, mettere un nome al token, selezionare tutte le caselle e generare il token, copiare il token e segnarsi il token da qualche parte.
