// Mensagem de erro da chamada de teste de notificação (429 tem texto próprio).
export function testErrorText(e) {
  return e.status === 429
    ? "Muitos testes seguidos. Tente de novo em 1 minuto."
    : e.message;
}
