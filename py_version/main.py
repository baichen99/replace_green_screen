import numpy as np
import cv2


def HexToRGB(hex):
    if hex[0] == '#':
        hex = hex[1:]
    if len(hex) != 6:
        return
    return int(hex[0:2], 16), int(hex[2:4], 16), int(hex[4:6], 16)

def RGBToHex(r, g, b):
    r = hex(r)[2:]
    g = hex(g)[2:]
    b = hex(b)[2:]
    return r + g + b


def create_image(size_x, size_y, hex):
    r, g, b = HexToRGB(hex)
    img = np.zeros([size_x, size_y, 3], np.uint8)
    img[:, :, 0] = np.ones([size_x, size_y])*b
    img[:, :, 1] = np.ones([size_x, size_y])*g
    img[:, :, 2] = np.ones([size_x, size_y])*r
    return img


def Convert(srcPath, dstPath, bg="#e6e6e6"):
    lower_green = np.array([68, 84, 153])
    upper_green = np.array([89, 255, 255])
    kernel = np.ones((1, 1), np.uint8)


    cap = cv2.VideoCapture(srcPath)
    fourcc = cv2.VideoWriter_fourcc(*"avc1")
    fps = int(cap.get(cv2.CAP_PROP_FPS))
    width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))
    height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))
    size = (width, height)
    # print(type(size))
    writer = cv2.VideoWriter(dstPath, fourcc, fps, size)

    while True:
        ok, frame = cap.read()
        if ok:
            # 获取mask
            hsv = cv2.cvtColor(frame, cv2.COLOR_BGR2HSV)
            mask = cv2.bitwise_not(cv2.inRange(hsv, lower_green, upper_green))
            # 腐蚀运算去绿边
            erosion = cv2.erode(mask, kernel)

            pure_gray = create_image(frame.shape[0], frame.shape[1], bg)

            # 获取人物
            person = cv2.bitwise_and(frame, frame, mask=mask)
            result = cv2.copyTo(person, erosion, pure_gray)

            writer.write(result)
        else:
            break

    cap.release()
    writer.release()


if __name__ == "__main__":
    Convert("video.mp4", "output.mp4", "#e6e6e6")