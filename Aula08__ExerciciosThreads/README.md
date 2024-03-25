1) Estudar os exemplos fornecidos na última aula (Aula 07) . O objetivo dos exemplos é apresentar o funcionamento da API da biblioteca pthreads, que permite desenvolver aplicações com múltiplos fios de execução, seguindo o modelo de threading estudado em aula.

Para executar os exemplos da Aula 07, baixe os mesmos do Moodle e abra um terminal (bash). Para verificar cada um dos exemplos, é necessário compilar e executar individualmente cada um, de acordo com os comandos a seguir. É importante que os códigos fonte sejam analisados e compreendidos.

```
gcc ex1.c -o ex1 -pthread
```

```
./ex1
```

2) Estudar e executar os exemplos fornecidos na aula de hoje (Aula 08). O objetivo dos exemplos é apresentar os problemas existentes com relação à sincronização de threads. Nas próximas aulas estudaremos as diferentes formas de coordenação entre threads/processos.

3) Você consegue responder as seguintes perguntas ou realizar os desafios sobre os exemplos da aula de hoje?
   1) Como as threads executam sempre na mesma sequência? O que aconteceria se as variáveis m1, m2 e m3 fossem inicializadas todas com o valor 1 no primeiro exemplo?
   - Isso ocorre devido ao semáforo implementado no código por meio de locks e unlocks nos mutexes. Eles ainda respeitam a mesma ordem e acabam em um ciclo, pelo mesmo motivo, porém, a sequência inicial é aleatória - não se sabe qual começa primeiro e tranca os seguintes.
   2) Sobre o mesmo exemplo, crie dois grupos de threads, onde as threads de um mesmo grupo possuem sequencia entre si, mas o mesmo não é necessário entre os dois grupos (as threads são independentes dois dois grupos).
   - Ver arquivo `ex0_grupos.c`
   3) Tente explicar o que acontece no segundo exemplo. Por que a classe resource_s, que implementa uma variável e métodos para ler e escrever nessa variável, não implementa sincronização (acesso atômico aos dados)? Por que isso precisa ser feito no contexto das threads? Como mutexes podem ser utilizados para resolver isso?
   - No segundo exemplo, 3 processos acessam a mesma variável e tentam incrementa-la. Diversas vezes, porém, a variável é incrementada mais de uma vez para o mesmo valor, o que não deveria acontecer. Isso ocorre pois a variável é acessada por mais de uma thread ao mesmo tempo, o que pode causar problemas de concorrência. A classe `resource_s` não implementa sincronização pois não possui locks e unlocks nos métodos de incremento e decremento. Isso precisa ser feito no contexto das threads pois são elas que acessam a variável de forma concorrente. Mutexes podem ser utilizados para resolver isso, pois garantem que apenas uma thread por vez acesse a variável.
