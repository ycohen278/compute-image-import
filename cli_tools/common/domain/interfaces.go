//  Copyright 2019 Google Inc. All Rights Reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package domain

import (
	"context"
	"io"
	"net/http"

	"cloud.google.com/go/storage"
)

// StorageClientInterface represents GCS storage client
type StorageClientInterface interface {
	CreateBucket(bucketName string, project string, attrs *storage.BucketAttrs) error
	UpdateBucket(bucketName string, attrs storage.BucketAttrsToUpdate) error
	Buckets(projectID string) *storage.BucketIterator
	GetBucketAttrs(bucket string) (*storage.BucketAttrs, error)
	GetBucket(bucket string) *storage.BucketHandle
	GetObject(bucket string, objectPath string) StorageObject
	GetObjects(bucket string, objectPath string) ObjectIteratorInterface
	GetObjectAttrs(bucket string, objectPath string) (*storage.ObjectAttrs, error)
	FindGcsFile(gcsDirectoryPath string, fileExtension string) (*storage.ObjectHandle, error)
	FindGcsFileDepthLimited(gcsDirectoryPath string, fileExtension string, lookupDepth int) (*storage.ObjectHandle, error)
	GetGcsFileContent(gcsObject *storage.ObjectHandle) ([]byte, error)
	WriteToGCS(destinationBucketName string, destinationObjectPath string, reader io.Reader) error
	DeleteGcsPath(gcsPath string) error
	DeleteObject(gcsPath string) error
	Close() error
}

// BucketIteratorCreatorInterface represents GCS bucket creator
type BucketIteratorCreatorInterface interface {
	CreateBucketIterator(ctx context.Context, storageClient StorageClientInterface,
		projectID string) BucketIteratorInterface
}

// BucketIteratorInterface represents GCS bucket iterator
type BucketIteratorInterface interface {
	Next() (*storage.BucketAttrs, error)
}

// ObjectIteratorCreatorInterface represents GCS object iterator creator
type ObjectIteratorCreatorInterface interface {
	CreateObjectIterator(bucket string, objectPath string) ObjectIteratorInterface
}

// ObjectIteratorInterface represents GCS Object iterator
type ObjectIteratorInterface interface {
	Next() (*storage.ObjectAttrs, error)
}

// TarGcsExtractorInterface represents TAR GCS extractor responsible for extracting TAR archives from GCS to
// GCS
type TarGcsExtractorInterface interface {
	ExtractTarToGcs(tarGcsPath string, destinationGcsPath string) error
}

// StorageObjectCreatorInterface represents GCS object creator
type StorageObjectCreatorInterface interface {
	GetObject(bucket string, objectPath string) StorageObject
}

// StorageObject represents GCS Object
type StorageObject interface {
	Delete() error
	GetObjectHandle() *storage.ObjectHandle
	NewReader() (io.ReadCloser, error)
	NewWriter() io.WriteCloser
	ObjectName() string
	Compose(src ...StorageObject) (*storage.ObjectAttrs, error)
	CopyFrom(src StorageObject) (*storage.ObjectAttrs, error)
}

// MetadataGCEInterface represents GCE metadata
type MetadataGCEInterface interface {
	OnGCE() bool
	Zone() (string, error)
	ProjectID() (string, error)
}

// ZoneValidatorInterface represents zone validator which validates if a zone is valid in given
// project
type ZoneValidatorInterface interface {
	ZoneValid(project string, zone string) error
}

// ScratchBucketCreatorInterface represents Daisy scratch (temporary) bucket creator
// To rebuild the mock, run `go generate ./...`
//
//go:generate go run github.com/golang/mock/mockgen -package mocks -destination mocks/mock_scratch_bucket_creator.go github.com/GoogleCloudPlatform/compute-image-import/cli_tools/common/domain ScratchBucketCreatorInterface
type ScratchBucketCreatorInterface interface {
	CreateScratchBucket(sourceFileFlag string, projectFlag string, fallbackZone string,
		enableUniformBucketLevelAccess bool) (string, string, error)
	IsBucketInProject(project string, bucketName string) bool
}

// ResourceLocationRetrieverInterface represents Daisy GCE/GCS resource location retriever
type ResourceLocationRetrieverInterface interface {
	GetZone(storageRegion string, project string) (string, error)
	GetLargestStorageLocation(storageLocation string) string
}

// HTTPClientInterface represents HTTP client
type HTTPClientInterface interface {
	Get(url string) (resp *http.Response, err error)
}

// ResourceDeleter checks whether images exist. If so, it deletes them.
type ResourceDeleter interface {
	DeleteImagesIfExist(images []Image)
	DeleteDisksIfExist(disks []Disk)
}

// Image holds the project, name, and URI of a GCP disk image.
type Image interface {
	GetProject() string
	GetImageName() string
	GetURI() string
}

// Disk holds the project, name, zone and URI of a PD.
type Disk interface {
	GetProject() string
	GetDiskName() string
	GetZone() string
	GetURI() string
}
