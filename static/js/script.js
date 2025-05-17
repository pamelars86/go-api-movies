// Script para manejar errores de carga de imágenes
document.addEventListener('DOMContentLoaded', function() {
    // Detectar errores de carga de imágenes y mostrar imagen de reemplazo
    const imgElements = document.querySelectorAll('img');
    const fallbackUrl = '/static/img/no-poster.svg';
    
    // Precargar la imagen de respaldo para evitar un nuevo fallo
    const preloadFallback = new Image();
    preloadFallback.src = fallbackUrl;
    
    // Configurar el manejador de errores para cada imagen
    imgElements.forEach(img => {
        // Guardar la URL original por si queremos reintentarla más tarde
        const originalSrc = img.src;
        
        // Ya configurado con el atributo onerror en HTML, pero lo reforzamos aquí
        img.addEventListener('error', function(e) {
            // Prevenir bucles infinitos: comprobar que no estamos ya usando la imagen de respaldo
            if (this.src !== fallbackUrl && !this.dataset.fallbackApplied) {
                console.log('Error cargando imagen:', originalSrc, '- Usando respaldo');
                this.src = fallbackUrl;
                this.dataset.fallbackApplied = 'true'; // Marcar como procesada
                
                // Añadir clase para estilos específicos
                this.classList.add('fallback-image');
            }
        });
        
        // Si la imagen ya está rota cuando se carga la página
        if (img.complete && (img.naturalWidth === 0 || img.naturalHeight === 0)) {
            img.src = fallbackUrl;
            img.dataset.fallbackApplied = 'true';
            img.classList.add('fallback-image');
        }
    });

    // Activar tooltips de Bootstrap
    const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
    tooltipTriggerList.map(function (tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl);
    });
}); 