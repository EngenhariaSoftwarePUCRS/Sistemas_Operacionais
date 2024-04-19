#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/wait.h>
#include <time.h>

char getRandRPS();
pid_t getWinner(char choice1, char choice2);
void printWinner(pid_t winner);

int main() {
    pid_t pK, winner;
    int ch1[2], ch2[2];
    pipe(ch1); pipe(ch2);

    // 'R' = Rock, 'P' = Paper, 'S' = Scissors
    char choice, choicePlayer1;

    pK = fork();
    if (pK < 0) { exit(1); }

    if (pK == 0) {
        close(ch1[1]);
        close(ch2[0]);
        choice = getRandRPS();
        printf("Escolha do processo %d (filho): %c\n", getpid(), choice);
        read(ch1[0], &choicePlayer1, sizeof(char));
        close(ch1[0]);
        winner = getWinner(choice, choicePlayer1);
        write(ch2[1], &winner, sizeof(pid_t));
        close(ch2[1]);
        exit(0);
    }

    if (pK > 0) {
        close(ch1[0]);
        close(ch2[1]);
        choice = getRandRPS();
        printf("Escolha do processo %d (pai): %c\n", getpid(), choice);
        write(ch1[1], &choice, sizeof(char));
        close(ch1[1]);
        read(ch2[0], &winner, sizeof(pid_t));
        close(ch2[0]);
        printWinner(winner);
        waitpid(pK, NULL, 0);
        exit(0);
    }

    printf("Obrigado por jogar!");

    return 0;
}

char getRandRPS() {
    srand(getpid());
    char rps[3] = {'R', 'P', 'S'};
    return rps[rand() % 3];
}

pid_t getWinner(char choice1, char choice2) {
    if (choice1 == choice2) return 0;
    if (choice1 == 'R' && choice2 == 'S') return getpid();
    if (choice1 == 'S' && choice2 == 'P') return getpid();
    if (choice1 == 'P' && choice2 == 'R') return getpid();
    return 0;
}

void printWinner(pid_t winner) {
    if (winner == 0) {
        printf("Empate!\n");
        return;
    }
    if (winner == getpid()) {
        printf("Vitória!\n");
        return;
    }
    printf("Derrota!\n");
}
