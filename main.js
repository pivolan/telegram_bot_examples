// Загружаем модели распознавания лиц из Face-API.js
Promise.all([
    faceapi.nets.faceRecognitionNet.loadFromUri('/models'),
    faceapi.nets.faceLandmark68Net.loadFromUri('/models'),
    faceapi.nets.ssdMobilenetv1.loadFromUri('/models')
]).then(startFaceRecognition);

// Начинаем распознавание лиц
function startFaceRecognition() {
    // Загружаем изображение для распознавания
    var image = document.createElement('img');
    image.src = 'path/to/image.jpg';
    image.addEventListener('load', function() {
        // Получаем контекст рисования для элемента <canvas>
        var canvas = document.createElement('canvas');
        canvas.width = image.width;
        canvas.height = image.height;
        var context = canvas.getContext('2d');
        context.drawImage(image, 0, 0, canvas.width, canvas.height);

        // Распознаем лица на изображении
        faceapi.detectAllFaces(canvas).then(function(detections) {
            // Выводим количество обнаруженных лиц на консоль
            console.log('Обнаружено лиц: ' + detections.length);

            // Рисуем прямоугольники вокруг обнаруженных лиц
            detections.forEach(function(detection) {
                var box = detection.box;
                context.strokeStyle = 'green';
                context.lineWidth = 2;
                context.strokeRect(box.x, box.y, box.width, box.height);
            });

            // Выводим обработанное изображение на страницу
            document.body.appendChild(canvas);
        });
    });
}
