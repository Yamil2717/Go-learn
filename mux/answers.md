# ¿Alguna respuesta que te sorprenda?
Si.

En `/blog`, `/hello` y `/world` esperaba un 404 porque esos endpoints no están registrados en el mux, pero los tres devuelven **"This is the home page."**

## ¿Por qué pasa esto?

En Go, `http.NewServeMux()` hace la coincidencia de rutas de una forma particular: no exige que la URL sea exactamente igual al patrón registrado. Un patrón que termina en `/`, como `"/"`, funciona como una raíz catch-all: cualquier URL de la que sea prefijo cae ahí.

Entonces, cuando llega una petición a `/blog`, el mux busca la ruta registrado más larga que coincida. Ni `/about` ni otro patrón más específico coinciden, así que termina en la ruta `/`, que corre `homeHandler`.

Por eso la única ruta que respondió distinto fue `/about`, porque sí tiene su propio handler registrado. Para que `/blog`, `/hello` y `/world` respondieran 404 habría que registrarlos con sus propios handlers, o que el handler de la raíz verifique la URL exacta.
