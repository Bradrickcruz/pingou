// Mostra a versão clara ou escura de uma imagem (ex.: logo) só com CSS.
// As duas ficam no HTML e o tema (data-theme no <html>) esconde uma, então a
// troca já vale na primeira pintura e não depende do React saber o tema.
// A escondida fica com display:none, fora da árvore de acessibilidade.
export function ThemedImage({ light, dark, className = "", ...rest }) {
  return (
    <>
      <img src={light} className={`${className} dark:hidden`} {...rest} />
      <img src={dark} className={`${className} hidden dark:inline`} {...rest} />
    </>
  );
}
