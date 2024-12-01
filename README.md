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
Dopo aver eseguito il comando non deve dare errori/scritte, se da delle scritte non sono state cambiate tutte le occorrenze del modulo nei file delle cartelle specificate, ricontrollare i file e rieseguire il comando vendor. 

Se avete già caricato l'api su github, andate sulla repo github, clonate (se non lo avete già fatto) la vostra repo, dopo averla clonata spostate il contenuto della cartella fantastic-coffee-decaffeinated (tranne la cartella doc) nella vostra cartella appena clonata, dopo aver fatto questo passaggio, eseguite il comando:
```bash
git add .
git commit -m "First commit"
git push
```

Se chiede l'accesso, mettere username e come password mettete un token generato da github. Per generare il token andate su github, cliccate sulla vostra immagine in alto a destra, cliccate su settings, andate su Developer settings, cliccate su personal access tokens, token classi, generate token classic, mettere un nome al token, selezionare tutte le caselle e generare il token, copiare il token e segnarsi il token da qualche parte. (Tranne nella chat con un utente in quanto per aprire una conversazione con un altro utente bisogna mandare un messaggio)

## Controllo dei file go

Per controllare la sintassi dei file go come il prof, bisogna scaricare golangci-lint, per farlo eseguite da terminale il seguente comando:
```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.62.0
```
Dopo averlo eseguito controllare se funziona scrivendo sempre da terminale `golangci-lint --version`, se vi esce la versione allora è stato installato correttamente, altrimenti provate a riavviare il terminale e a controllare di nuovo la versione.
Per controllare i file, eseguito questi due comandi 
```bash
golangci-lint run service/api -E rowserrcheck
```
Questo comando server per controllare i file all'interno della cartella api, per controllare i file all'interno della cartella database eseguire il comando:
```bash
golangci-lint run service/database -E rowserrcheck
``` 

## Test API con REST API CLIENT (Estensione VS CODE)

Scaricare l'estensione REST API CLIENT, questo è il link dell'esentione se non riuscite a trovarlo:
- https://marketplace.visualstudio.com/items?itemName=donebd.rest-client-api

Dopo averlo scaricato create un file .http all'interno della cartella service/api, come per esempio test.https. All'interno del file 
dovete scrivere le richieste http da mandare al proprio server. Per vedere come scrivere i file aprire il mio file chiamato ciccio.http all'interno della cartella service/api.
Dopo aver creato il file, nel terminale eseguite il seguente comando 
```bash
go run ./cmd/webapi
```
Dopo averlo eseguito su vs code andate nel file .http che avete creato e se l'estensione funziona sopra ogni richiesta appare una 
piccola scritta Send Request, cliccate sulla scritta e vi appare una finestra nella quale viene riportata il messaggio di risposta alla richiesta http come lo avete impostato voi nei.

N.B. Quando eseguito il comando go run vi potrebbero dare degli errori, conrtollate bene gli errori e risolveteli prima di passare all'estensione. Inoltre, se il comando è andato a buon fine vi spunta la seguente scritta:
```bash
INFO[0000] API Listening on 0.0.0.0:3000
```
Se il numero non è 3000 va bene lo stesso, ma all'interno del file http assicuratevi che il link http abbia il numero che vi spunta al posto di 3000 (Guardate il mio file .http per capire)
