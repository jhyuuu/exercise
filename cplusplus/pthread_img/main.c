#include <pthread.h>

#include <vector>
#include <thread>
#include <string>
#include <fstream>

#include <opencv2/opencv.hpp>
using namespace cv;

struct CVImg {
  int index;
  std::string img_path;
  Mat mat;
};

std::vector<std::string> get_imgs(const std::string& img_list) {
  std::vector<std::string> imgs;
  
  std::string line;
  std::ifstream in(img_list);
  if (in.is_open()) {
    while (!in.eof()) {
      getline(in, line);
      if (line.size() != 0)
        imgs.push_back(line);
    }
  }

  return imgs;
}

void load_image(int index, const std::string& img, std::vector<CVImg>* cv_imgs) {
  Mat mat = imread(img);
  CVImg cv_img = {index, img, mat};
  cv_imgs->push_back(cv_img);

  printf("tt %d: %s\n", index, img.c_str());
}

void load_images(const std::vector<std::string>& imgs) {
  int thrds = 8;
  std::thread threads[thrds];
  std::vector<CVImg> cv_imgs;
  
  int batches = imgs.size() / thrds;
  int left_over = imgs.size() % thrds;

  for (int b = 0; b < batches; ++b) {
    for (int i = 0; i < thrds; ++i) {
      int index = b*thrds + i;
      std::string img = imgs[index];
      printf("aa %d: %s\n", index, img.c_str());
      threads[i] = std::thread(load_image, index, img, &cv_imgs);
    }
    for (int i = 0; i < thrds; ++i) {
      threads[i].join();
    }
    for (int i = 0; i < cv_imgs.size(); ++i) {
      char buf[512];
      CVImg cv_img = cv_imgs[i];
      sprintf(buf, "%d %s", cv_img.index, cv_img.img_path.c_str());
      imshow(buf, cv_img.mat);
      waitKey(0);// 延时时间
    }
    printf("imgs: %lu\n", cv_imgs.size());
    cv_imgs.clear();
  }

  if (left_over > 0) {
    for (int i = 0; i < left_over; ++i) {
      int index = batches*thrds + i;
      std::string img = imgs[index];
      printf("aa %d: %s\n", index, img.c_str());
      threads[i] = std::thread(load_image, index, img, &cv_imgs);
    }
    for (int i = 0; i < left_over; ++i) {
      threads[i].join();
    }
    for (int i = 0; i < cv_imgs.size(); ++i) {
      char buf[512];
      CVImg cv_img = cv_imgs[i];
      sprintf(buf, "%d %s", cv_img.index, cv_img.img_path.c_str());
      imshow(buf, cv_img.mat);
      waitKey(0);// 延时时间
    }
    printf("imgs: %lu\n", cv_imgs.size());
    cv_imgs.clear();
  }
}

std::thread load_data(const std::vector<std::string>& imgs) {
  std::thread thread(load_images, imgs);
  
  return thread;
}

int main() {
  std::string img_list = "img_list.txt";
  std::vector<std::string> imgs = get_imgs(img_list);
  std::thread thread = load_data(imgs);
  thread.join();
  printf("after thread in load_data\n");
  return 0;
}
