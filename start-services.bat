start powershell -NoExit -Command "$env:PORT='3000'; npm run start:dev:producer"
start powershell -NoExit -Command "$env:PORT='3001'; npm run start:dev:consumer"
start powershell -NoExit -Command "$env:PORT='3002'; npm run start:dev:notification"